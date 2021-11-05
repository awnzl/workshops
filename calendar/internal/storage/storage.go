package storage

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // driver initialization
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresQL(l *log.Logger) *PostgresDB {
	return &PostgresDB{}
}

func (s *PostgresDB) Connect(dbURL string) (err error) {
	s.db, err = sql.Open("postgres", dbURL)

	return
}

func (s *PostgresDB) Close() {
	s.db.Close()
}
