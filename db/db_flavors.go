package db

import "fmt"

func getCreateMetadataQueryForDriver(driver string) string {
	switch driver {
	case "mysql":
		return mySQLMetadataTableExists
	default: // postgres
		return  postgresMetadataTableExists
	}
}

func getInsertIntoMetadataQueryForDriver(driver string) string {
	switch driver {
	case "mysql":
		return mySQLInsertIntoMetadataTable
	default: // postgres
		return  postgresInsertIntoMetadataTable
	}
}

func getDeleteFromMetadataQueryForDriver(driver string) string {
	switch driver {
	case "mysql":
		return mySQLDeleteFromMetadataTable
	default: // postgres
		return  postgresDeleteFromMetadataTable
	}
}

func getSelectIDNameHashFromMetadataWhereNameQueryForDriver(driver string) string {
	switch driver {
	case "mysql":
		return mySQLSelectIDNameHashFromMetadataTableWhereName
	default: // postgres
		return  postgresSelectIDNameHashFromMetadataTableWhereName
	}
}

func getSelectIDNameHashFromMetadataWhereHashQueryForDriver(driver string) string {
	switch driver {
	case "mysql":
		return mySQLSelectIDNameHashFromMetadataTableWhereHash
	default: // postgres
		return  postgresSelectIDNameHashFromMetadataTableWhereHash
	}
}

func getConnectionInfoForDatabase(d Database) string {
	switch d.Driver {
	case "mysql":
		return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", d.User, d.Password, d.Host, d.Port, d.Dbname)
	default: // postgres
		return  fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",d.Host, d.Port, d.User, d.Password, d.Dbname)
	}
}