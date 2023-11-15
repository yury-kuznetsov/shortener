package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"math/rand"
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

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS urls (" +
		"code varchar not null constraint urls_pk unique," +
		"uri varchar not null constraint urls_pk2 unique" +
		")")
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
	db, err := sql.Open("pgx", DSN)
	if err != nil {
		return "", err
	}
	defer db.Close()

	key := generateKey()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err = db.ExecContext(ctx, "INSERT INTO urls (code, uri) VALUES($1,$2)", key, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := db.QueryRowContext(ctx, "SELECT code FROM urls WHERE uri = $1", value)
			if errScan := row.Scan(&key); errScan != nil {
				return "", errScan
			}
			return key, err
		}
		return "", err
	}

	return key, nil
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

func generateKey() string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var length = 8
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	key := make([]byte, length)
	for i := range key {
		key[i] = charset[seededRand.Intn(len(charset))]
	}

	return string(key)
}
