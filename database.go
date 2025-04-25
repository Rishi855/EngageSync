package main

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres dbname=engegesync sslmode=disable password=root")
	if err != nil {
		log.Fatal(err)
	}
}

func execQuery(query string, args ...interface{}) (sql.Result, error) {
	return db.Exec(query, args...)
}

func queryRow(query string, args ...interface{}) *sql.Row {
	return db.QueryRow(query, args...)
}

func queryRows(query string, args ...interface{}) (*sql.Rows, error) {
	return db.Query(query, args...)
}
