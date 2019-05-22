package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"strings"
)

type Postgres struct {
	Host string
	Port int
	User string
	Password string
	Dbname string
}

func (p Postgres) connect() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	ensureDeckardTableExists(*db)

	return db
}

func ensureDeckardTableExists(db sql.DB) (sql.Result, error) {
	sqlStatement := `
CREATE TABLE IF NOT EXISTS deckard_horadric_cube (
  id SERIAL,
  name TEXT NOT NULL,
  hash TEXT NOT NULL
)`

	result, err := db.Exec(sqlStatement)

	if err != nil {
		log.Fatal("Failed to ensure that the metadata table used by Deckard exists. I tried to execute the following query and failed:", sqlStatement, err)
	}

	return result, err
}

func createHash(s string) string {
	hash := md5.New()
	_, err := io.WriteString(hash, s)
	if err != nil {
		log.Fatal("Failed to hash:", s, "\nTerminating...")
	}
	return hex.EncodeToString(hash.Sum(nil)[:])
}

func storeMigrationMetadata(db *sql.DB, query Query) (sql.Result, error) {
	sqlStatement := `INSERT INTO deckard_horadric_cube (name, hash) VALUES ($1, $2)`
	hash := createHash(query.Value)
	return db.Exec(sqlStatement, query.Name, hash)
}

func (p Postgres) RunUp(migration Migration) {
	db := p.connect()
	defer db.Close()
	ranSomething := false
	for _, query := range migration.Queries {
		if hasDbAlreadyRan(db, query) {
			continue
		}
		println(query.Value)
		_, err := db.Exec(query.Value)
		ranSomething = true
		if err == nil {
			_, err = storeMigrationMetadata(db, query)
			if err != nil {
				log.Fatal("Failed to write migration metadata.\n", err)
			}
		} else {
			log.Fatal("Failed to execute migration!", err)
		}
	}

	if !ranSomething {
		fmt.Println("No migrations were ran!")
	}
}

func deleteMigrationMetadata(db *sql.DB, query Query) (sql.Result, error) {
	sqlStatement := `DELETE FROM deckard_horadric_cube WHERE name = $1`
	upName := strings.ReplaceAll(query.Name, ".down.sql", ".up.sql")
	return db.Exec(sqlStatement, upName)
}

func (p Postgres) RunDown(migration Migration) {
	db := p.connect()
	defer db.Close()
	ranSomething := false
	for _, query := range migration.Queries {
		if !canRunDownMigration(db, query) {
			continue
		}
		println(query.Value)
		_, err := db.Exec(query.Value)
		ranSomething = true
		if err == nil {
			_, err = deleteMigrationMetadata(db, query)
			if err != nil {
				log.Fatal("Failed to delete migration metadata.")
			}
		} else {
			log.Fatal("Failed to execute migration!", err)
		}
	}

	if !ranSomething {
		fmt.Println("No migrations were ran!")
	}
}

func canRunDownMigration(db *sql.DB, query Query) bool {
	var name string
	var hash string
	var id int
	sqlStatement := `SELECT id, name, hash FROM deckard_horadric_cube WHERE name = $1`
	upName := strings.ReplaceAll(query.Name, ".down.sql", ".up.sql")
	row := db.QueryRow(sqlStatement, upName)
	switch err := row.Scan(&id, &name, &hash); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		return false
	}
}

func hasDbAlreadyRan(db *sql.DB, query Query) bool {
	var name string
	var hash string
	var id int
	sqlStatement := `SELECT id, name, hash FROM deckard_horadric_cube WHERE hash=$1;`
	row := db.QueryRow(sqlStatement, createHash(query.Value))
	switch err := row.Scan(&id, &name, &hash); err {
	case sql.ErrNoRows:
		return false
	case nil:
		return true
	default:
		return false
	}
}

func (p Postgres) Verify(migration Migration) {
	db := p.connect()
	defer db.Close()
	for _, query := range migration.Queries {
		println("Verifying:", query.Value)

		if hasDbAlreadyRan(db, query) {
			fmt.Println(`Validation Successful! It looks like you've already ran`, query.Value, `on this database.`)
		} else {
			fmt.Println(`Warning: Deckard cannot verify the migration.
Please ensure that the migration has not been changed locally since it was last ran.
If the migration has been changed, you may want to run deckard down and deckard up again.
Consider backing up your data before running deckard down.`)
		}
	}
}