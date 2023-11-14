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

	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return &Storage{}, err
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (code varchar not null constraint urls_pk unique, uri varchar not null)")
	if err != nil {
		return &Storage{}, err
	}

	return &Storage{}, nil
}

func (s *Storage) Get(code string) (string, error) {
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return "", err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := db.QueryRowContext(ctx, "SELECT uri FROM urls WHERE code = $1", code)

	var uri string
	if err = row.Scan(&uri); err != nil {
		return "", err
	}

	return uri, nil
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
