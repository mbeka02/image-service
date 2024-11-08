package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type Store struct {
	*Queries
}

func NewStore(uri string) (*Store, error) {
	conn, err := sql.Open("postgres", uri)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return &Store{
		Queries: New(conn),
	}, err
}
