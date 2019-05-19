package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type postgres struct {
	host string
	port int
	user string
	password string
	dbname string
}

func (p postgres) connect() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		p.host, p.port, p.user, p.password, p.dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}

func (p postgres) runUp() {

}

func (p postgres) runDown() {

}