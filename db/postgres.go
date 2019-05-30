package db

import (
	_ "github.com/lib/pq"
)

const postgresMetadataTableExists string = `
CREATE TABLE IF NOT EXISTS deckard_horadric_cube (
  id SERIAL,
  name TEXT NOT NULL,
  hash TEXT NOT NULL
)`

const postgresInsertIntoMetadataTable string =
	`INSERT INTO deckard_horadric_cube (name, hash) VALUES ($1, $2)`

const postgresDeleteFromMetadataTable string =
	`DELETE FROM deckard_horadric_cube WHERE name = $1`

const postgresSelectIdNameHashFromMetadataTableWhereName string =
	`SELECT id, name, hash FROM deckard_horadric_cube WHERE name = $1`

const postgresSelectIdNameHashFromMetadataTableWhereHash string =
	`SELECT id, name, hash FROM deckard_horadric_cube WHERE hash=$1;`