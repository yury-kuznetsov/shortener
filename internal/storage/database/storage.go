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

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	s := Storage{db: db}

	if err != nil {
		return &s, err
	}

	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS urls (" +
		"code varchar not null constraint urls_pk unique," +
		"uri varchar not null constraint urls_pk2 unique" +
		")")

	return &s, err
}

func (s *Storage) Get(ctx context.Context, code string) (string, error) {
	//defer s.db.Close()
	row := s.db.QueryRowContext(ctx, "SELECT uri FROM urls WHERE code = $1", code)

	var uri string
	if err := row.Scan(&uri); err != nil {
		return "", err
	}

	return uri, nil
}

func (s *Storage) Set(ctx context.Context, value string) (string, error) {
	//defer s.db.Close()
	key := generateKey()

	_, err := s.db.ExecContext(ctx, "INSERT INTO urls (code, uri) VALUES($1,$2)", key, value)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			row := s.db.QueryRowContext(ctx, "SELECT code FROM urls WHERE uri = $1", value)
			if errScan := row.Scan(&key); errScan != nil {
				return "", errScan
			}
			return key, err
		}
		return "", err
	}

	return key, nil
}

func (s *Storage) HealthCheck(ctx context.Context) error {
	//defer s.db.Close()
	return s.db.PingContext(ctx)
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
