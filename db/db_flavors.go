package db

func getCreateMetadataQueryForDriver(driver string) string {
	return postgresMetadataTableExists
}

func getInsertIntoMetadataQueryForDriver(driver string) string {
	return postgresInsertIntoMetadataTable
}

func getDeleteFromMetadataQueryForDriver(driver string) string {
	return postgresDeleteFromMetadataTable
}

func getSelectIdNameHashFromMetadataWhereNameQueryForDriver(driver string) string {
	return postgresSelectIdNameHashFromMetadataTableWhereName
}

func getSelectIdNameHashFromMetadataWhereHashQueryForDriver(driver string) string {
	return postgresSelectIdNameHashFromMetadataTableWhereHash
}