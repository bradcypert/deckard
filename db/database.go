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

type Database struct {
	Host string
	Port int
	User string
	Password string
	Dbname string
}

func (d Database) connect(driver string) *sql.DB {
	connectionInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.Dbname)

	db, err := sql.Open(driver, connectionInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	ensureDeckardTableExists(driver, *db)

	return db
}

func ensureDeckardTableExists(driver string, db sql.DB) (sql.Result, error) {
	sqlStatement := getCreateMetadataQueryForDriver(driver)

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

func storeMigrationMetadata(driver string, db *sql.DB, query Query) (sql.Result, error) {
	sqlStatement := getInsertIntoMetadataQueryForDriver(driver)
	hash := createHash(query.Value)
	return db.Exec(sqlStatement, query.Name, hash)
}

func (d Database) RunUp(migration Migration) {
	const driver string = "postgres" // TODO: Change me
	db := d.connect(driver)
	defer db.Close()
	ranSomething := false
	for _, query := range migration.Queries {
		if hasDbAlreadyRan(driver, db, query) {
			continue
		}
		println(query.Value)
		_, err := db.Exec(query.Value)
		ranSomething = true
		if err == nil {
			_, err = storeMigrationMetadata(driver, db, query)
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

func deleteMigrationMetadata(driver string, db *sql.DB, query Query) (sql.Result, error) {
	sqlStatement := getDeleteFromMetadataQueryForDriver(driver)
	upName := strings.ReplaceAll(query.Name, ".down.sql", ".up.sql")
	return db.Exec(sqlStatement, upName)
}

func (p Database) RunDown(migration Migration) {
	const driver string = "postgres"
	db := p.connect(driver)
	defer db.Close()
	ranSomething := false
	for _, query := range migration.Queries {
		if !canRunDownMigration(driver, db, query) {
			continue
		}
		println(query.Value)
		_, err := db.Exec(query.Value)
		ranSomething = true
		if err == nil {
			_, err = deleteMigrationMetadata(driver, db, query)
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

func canRunDownMigration(driver string, db *sql.DB, query Query) bool {
	var name string
	var hash string
	var id int
	sqlStatement := getSelectIdNameHashFromMetadataWhereNameQueryForDriver(driver)
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

func hasDbAlreadyRan(driver string, db *sql.DB, query Query) bool {
	var name string
	var hash string
	var id int
	sqlStatement := getSelectIdNameHashFromMetadataWhereHashQueryForDriver(driver)
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

func (d Database) Verify(migration Migration) {
	const driver string = "postgres"
	db := d.connect(driver)
	defer db.Close()
	for _, query := range migration.Queries {
		println("Verifying:", query.Value)

		if hasDbAlreadyRan(driver, db, query) {
			fmt.Println(`Validation Successful! It looks like you've already ran`, query.Value, `on this database.`)
		} else {
			fmt.Println(`Warning: Deckard cannot verify the migration.
Please ensure that the migration has not been changed locally since it was last ran.
If the migration has been changed, you may want to run deckard down and deckard up again.
Consider backing up your data before running deckard down.`)
		}
	}
}