package db

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"crypto/md5"
	_ "github.com/lib/pq"
	"io"
	"log"
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

	fmt.Println("Successfully connected!")

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

	result, error := db.Exec(sqlStatement)

	if error != nil {
		log.Fatal("Failed to ensure that the metadata table used by Deckard exists. I tried to execute the following query and failed:", sqlStatement, error)
	}

	return result, error
}

func storeMigrationMetadata(db *sql.DB, query Query) (sql.Result, error) {
	sqlStatement := `INSERT INTO deckard_horadric_cube (name, hash) VALUES ($1, $2)`
	hash := md5.New()
	_, err := io.WriteString(hash, query.Value)
	if err != nil {
		log.Fatal("Failed to hash query:", query.Value, "\nTerminating...")
	}
	return db.Exec(sqlStatement, query.Name,  hex.EncodeToString(hash.Sum(nil)[:]))
}

func (p Postgres) RunUp(migration Migration) {
	db := p.connect()
	defer db.Close()
	println("Running...")
	for _, query := range migration.Queries {
		println(query.Value)
		_, err := db.Exec(query.Value)
		if err != nil {
			_, err = storeMigrationMetadata(db, query)
			if err != nil {
				log.Fatal("Failed to write migration metadata.\n", err)
			}
		}
	}
}

func deleteMigrationMetadata(db *sql.DB, query Query) (sql.Result, error) {
	sqlStatement := `DELETE FROM deckard_horadric_cube WHERE name = $1 AND hash = $2`
	hash := md5.New()
	_, err := io.WriteString(hash, query.Value)
	if err != nil {
		log.Fatal("Failed to hash query:", query.Value, "\nTerminating...")
	}
	return db.Exec(sqlStatement, query.Name,  hex.EncodeToString(hash.Sum(nil)[:]))
}

func (p Postgres) RunDown(migration Migration) {
	db := p.connect()
	defer db.Close()
	for _, query := range migration.Queries {
		_, err := db.Exec(query.Value)
		if err != nil {
			_, err = deleteMigrationMetadata(db, query)
			if err != nil {
				log.Fatal("Failed to write migration metadata.")
			}
		}
	}
}