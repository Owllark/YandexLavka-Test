package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type PostgreSQLDatabase struct {
	conn *sql.DB
	PostgreSQLQuery
}

func (p *PostgreSQLDatabase) Connect() error {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=lavka host=localhost sslmode=disable")
	if err != nil {
		return err
	}
	p.conn = db
	p.queryConn = db
	return nil
}

func (p *PostgreSQLDatabase) Close() error {
	return p.conn.Close()
}

type PostgreSQLQuery struct {
	queryConn *sql.DB
}

func (p *PostgreSQLQuery) Exec(query string, args ...interface{}) (sql.Result, error) {
	return p.queryConn.Exec(query, args...)
}

func (p *PostgreSQLQuery) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return p.queryConn.Query(query, args...)
}

func (p *PostgreSQLQuery) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.queryConn.QueryRow(query, args...)
}
