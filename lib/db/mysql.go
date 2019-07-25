package db

const mySQLMetadataTableExists string = `
CREATE TABLE IF NOT EXISTS deckard_horadric_cube (
  id int NOT NULL AUTO_INCREMENT,
  name TEXT NOT NULL,
  hash TEXT NOT NULL,
  primary key (id)
)`

const mySQLInsertIntoMetadataTable string = `INSERT INTO deckard_horadric_cube (name, hash) VALUES (?, ?)`

const mySQLDeleteFromMetadataTable string = `DELETE FROM deckard_horadric_cube WHERE name = ?`

const mySQLSelectIDNameHashFromMetadataTableWhereName string = `SELECT id, name, hash FROM deckard_horadric_cube WHERE name = ?`

const mySQLSelectIDNameHashFromMetadataTableWhereHash string = `SELECT id, name, hash FROM deckard_horadric_cube WHERE hash=?;`
