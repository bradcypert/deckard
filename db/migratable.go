package db

type migratable interface {
	runUp()
	runDown()
	connect()
}