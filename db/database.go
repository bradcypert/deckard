package db

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/lib/pq" // import for side effects
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"strings"
)

const ValidationError = `Warning: Deckard cannot verify the migration.
Please ensure that the migration has not been changed locally since it was last ran.
If the migration has been changed, you may want to run deckard down and deckard up again.
Consider backing up your data before running deckard down.`

const FailedToWriteMetadata = "Failed to write migration metadata.\n"
const FailedToDeleteMetadata = "Failed to delete migration metadata.\n"
const FailedToExecuteMigration = "Failed to execute migration!"
const NoMigrationsRan = "No Migrations Were Ran!"
const FailedToEnsureMetadataTable = "Failed to ensure that the metadata table used by Deckard exists. I tried to execute the following query and failed:"
const FailedToHash = "Failed to Hash:"

// Database a structure for defining a database connection string.
type Database struct {
	Host string
	Port int
	User string
	Password string
	Dbname string
	Driver string
}

// RunUp Runs an up migration against a given database.
func (d Database) RunUp(migration Migration) {
	db := d.connect(d.Driver)
	defer db.Close()
	ranSomething := false
	for _, query := range migration.Queries {
		if hasDbAlreadyRan(d.Driver, db, query) {
			continue
		}
		println(query.Value)
		_, err := db.Exec(query.Value)
		ranSomething = true
		if err == nil {
			_, err = storeMigrationMetadata(d.Driver, db, query)
			if err != nil {
				log.Fatal(FailedToWriteMetadata, err)
			}
		} else {
			log.Fatal(FailedToExecuteMigration, err)
		}
	}

	if !ranSomething {
		fmt.Println(NoMigrationsRan)
	}
}

// RunDown Runs a down migration against a given database.
func (d Database) RunDown(migration Migration) {
	db := d.connect(d.Driver)
	defer db.Close()
	ranSomething := false
	for _, query := range migration.Queries {
		if !canRunDownMigration(d.Driver, db, query) {
			continue
		}
		println(query.Value)
		_, err := db.Exec(query.Value)
		ranSomething = true
		if err == nil {
			_, err = deleteMigrationMetadata(d.Driver, db, query)
			if err != nil {
				log.Fatal(FailedToDeleteMetadata)
			}
		} else {
			log.Fatal(FailedToExecuteMigration, err)
		}
	}

	if !ranSomething {
		fmt.Println(NoMigrationsRan)
	}
}

// Verify verifies that a given migration has been ran against a given database.
func (d Database) Verify(migration Migration) {
	db := d.connect(d.Driver)
	defer db.Close()
	for _, query := range migration.Queries {
		println("Verifying:", query.Value)

		if hasDbAlreadyRan(d.Driver, db, query) {
			fmt.Println(`Validation Successful! It looks like you've already ran`, query.Value, `on this database.`)
		} else {
			fmt.Println(ValidationError)
		}
	}
}

func (d Database) connect(driver string) *sql.DB {
	connectionInfo := getConnectionInfoForDatabase(d)

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
		log.Fatal(FailedToEnsureMetadataTable, sqlStatement, err)
	}

	return result, err
}

func createHash(s string) string {
	hash := md5.New()
	_, err := io.WriteString(hash, s)
	if err != nil {
		log.Fatal(FailedToHash, s, "\nTerminating...")
	}
	return hex.EncodeToString(hash.Sum(nil)[:])
}

func storeMigrationMetadata(driver string, db *sql.DB, query Query) (sql.Result, error) {
	sqlStatement := getInsertIntoMetadataQueryForDriver(driver)
	hash := createHash(query.Value)
	return db.Exec(sqlStatement, query.Name, hash)
}

func deleteMigrationMetadata(driver string, db *sql.DB, query Query) (sql.Result, error) {
	sqlStatement := getDeleteFromMetadataQueryForDriver(driver)
	upName := strings.ReplaceAll(query.Name, ".down.sql", ".up.sql")
	return db.Exec(sqlStatement, upName)
}

func canRunDownMigration(driver string, db *sql.DB, query Query) bool {
	var name string
	var hash string
	var id int
	sqlStatement := getSelectIDNameHashFromMetadataWhereNameQueryForDriver(driver)
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
	sqlStatement := getSelectIDNameHashFromMetadataWhereHashQueryForDriver(driver)
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