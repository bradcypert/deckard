package db

type Migration struct {
	Queries []Query
}

type Query struct {
	Name  string
	Value string
}
