package database

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yury-kuznetsov/shortener/internal/models"
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
		"uri varchar not null constraint urls_pk2 unique," +
		"user_id integer default 0 not null" +
		")")

	return &s, err
}

func (s *Storage) Get(ctx context.Context, code string, userID int) (string, error) {
	//defer s.db.Close()
	row := s.db.QueryRowContext(
		ctx,
		"SELECT uri FROM urls WHERE code = $1 AND user_id = $2",
		code, userID,
	)

	var uri string
	if err := row.Scan(&uri); err != nil {
		return "", err
	}

	return uri, nil
}

func (s *Storage) Set(ctx context.Context, value string, userID int) (string, error) {
	//defer s.db.Close()
	key := generateKey()

	_, err := s.db.ExecContext(
		ctx,
		"INSERT INTO urls (code, uri, user_id) VALUES($1,$2,$3)",
		key, value, userID,
	)
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

func (s *Storage) GetByUser(ctx context.Context, userID int) ([]models.GetByUserResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT code, uri FROM urls WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}

	response := make([]models.GetByUserResponse, 0)

	for rows.Next() {
		var data models.GetByUserResponse
		if err = rows.Scan(&data.ShortURL, &data.OriginalURL); err != nil {
			return nil, err
		}
		response = append(response, data)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return response, nil
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
