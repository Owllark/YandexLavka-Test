package db

import "database/sql"

type Database interface {
	Connect() error
	Close() error
	Query
}

type Query interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}
