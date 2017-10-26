package testdb

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)


func Setup() *sql.DB {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		panic(err)
	}
	return db
}
