package repository

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func NewDBConfig() *sql.DB {

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	return db
}
