package db

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func GetDB() *sql.DB {
	//"postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	postgress_url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		"postgres", 5432, "postgres", "postgres", "postgres")

	db, err := sql.Open("postgres", postgress_url)

	if err != nil {
		panic(err.Error())
	}

	return db
}
