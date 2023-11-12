package database

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Storage struct{}

var DSN string

func NewStorage(dsn string) (*Storage, error) {
	DSN = dsn
	return &Storage{}, nil
}

func (s *Storage) Get(code string) (string, error) {
	return "", nil
}

func (s *Storage) Set(value string) (string, error) {
	return "", nil
}

func (s *Storage) HealthCheck() error {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
