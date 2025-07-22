package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

func Inint() {
	var err error

	DB, err = sql.Open("postgres", "host=localhost port=5432 user=postgres password=postgres sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal(err)
	}
}
