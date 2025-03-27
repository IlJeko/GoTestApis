package main

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// db config
func NewDB(dbDriver string, dbSource string) *sql.DB {
	db, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		panic(err)
	}

	return db
}
