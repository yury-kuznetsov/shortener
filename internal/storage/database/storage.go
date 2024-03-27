// Package database provides functionality for interacting with a database.
// DBConfig represents the configuration details for connecting to a database.
package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/yury-kuznetsov/shortener/internal/models"
)

// Storage is a type that represents a storage object
type Storage struct {
	db *sql.DB
}

// NewStorage creates a new instance of Storage initialized with a PostgreSQL database connection.
// It takes a DSN (Data Source Name) string as a parameter and returns a pointer to a Storage instance and an error.
func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("pgx", dsn)
	s := Storage{db: db}

	if err != nil {
		return &s, err
	}

	_, err = s.db.Exec("CREATE TABLE IF NOT EXISTS urls (" +
		"code varchar not null constraint urls_pk unique," +
		"uri varchar not null constraint urls_pk2 unique," +
		"user_id integer default 0 not null," +
		"is_deleted boolean default false not null" +
		")")

	return &s, err
}

// Get retrieves the original URI associated with the given code and user ID.
// It takes a context, a code string, and a userID int as parameters.
// It returns the URI string and an error.
// The code queries the database to fetch the URI and is_deleted flag for the given code.
// If the row scan fails, it returns an error.
// If the is_deleted flag is true, it returns models.ErrRowDeleted.
// Otherwise, it returns the URI string and nil error.
func (s *Storage) Get(ctx context.Context, code string, userID int) (string, error) {
	//defer s.db.Close()
	row := s.db.QueryRowContext(
		ctx,
		"SELECT uri, is_deleted FROM urls WHERE code = $1",
		code,
	)

	var uri string
	var isDeleted bool
	if err := row.Scan(&uri, &isDeleted); err != nil {
		return "", err
	}

	if isDeleted {
		return "", models.ErrRowDeleted
	}

	return uri, nil
}

// Set adds a new URL to the storage with a generated code.
//
// It takes a context, a value string, and a userID int as parameters.
// It returns the generated code string and an error.
//
// The method generates a new code using the `generateKey` function.
// It then attempts to insert the code, URI, and user ID into the `urls` table.
// If the insertion fails, it checks if the error is a unique violation error.
// If it is, it queries the `urls` table to find the existing code for the given URI.
// If the query and scan fail, it returns an empty string and the scan error.
// Otherwise, it returns the existing code and the original error.
//
// If the insertion is successful, it returns the generated code and nil error.
//
// Example usage:
//
//	value := "http://example.com"
//	userID := 123
//	code, err := storage.Set(ctx, value, userID)
//	if err != nil {
//	  log.Fatal(err)
//	}
//	fmt.Println("Generated code:", code)
//
// Note: The `Storage` type must have a field `db` of type `*sql.DB`.
// The `generateKey` function must be defined in the same package as the `Set` method.
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

// GetByUser retrieves all URLs associated with the given user ID.
// It takes a context and a userID int as parameters.
// It returns a slice of models.GetByUserResponse and an error.
// The code queries the database to fetch the code and URI for URLs associated with the given user ID.
// It scans the retrieved rows into models.GetByUserResponse struct and appends them to the response slice.
// If the row scan fails, it returns an error.
// It then checks for any errors during the iteration of the rows.
// If an error occurs, it returns the error.
// Otherwise, it returns the response slice and nil error.
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

// SoftDelete marks the URLs associated with the given messages as deleted.
// It takes a context and a slice of models.RmvUrlsMsg as parameters.
// It returns an error.
// The code iterates over the messages, constructing a WHERE clause for each message
// in the format "(code = $x AND user_id = $y)", where x is derived from the message's index and
// y is derived from the message's index plus one.
// It appends these params to a values slice and the message's code and userID to an args slice.
// It then constructs the UPDATE query with the WHERE clause joined by OR.
// Finally, it executes the query and returns any resulting error.
func (s *Storage) SoftDelete(ctx context.Context, messages []models.RmvUrlsMsg) error {
	var values []string
	var args []any

	for i, msg := range messages {
		base := i * 2
		params := fmt.Sprintf("(code = $%d AND user_id = $%d)", base+1, base+2)
		values = append(values, params)
		args = append(args, msg.Code, msg.UserID)
	}

	query := "UPDATE urls SET is_deleted = true WHERE " + strings.Join(values, " OR ") + ";"
	_, err := s.db.ExecContext(ctx, query, args...)

	return err
}

// HealthCheck performs a health check by pinging the underlying database.
// It takes a context as a parameter.
// It returns an error.
// The code uses the PingContext method of the db field to send a ping request to the database.
// It returns the error received from the PingContext method.
func (s *Storage) HealthCheck(ctx context.Context) error {
	//defer s.db.Close()
	return s.db.PingContext(ctx)
}

// GetStats retrieves the current statistics of the storage.
func (s *Storage) GetStats(ctx context.Context) (int, int, error) {
	row := s.db.QueryRowContext(ctx, "SELECT COUNT(DISTINCT user_id) FROM urls")
	var users int
	if err := row.Scan(&users); err != nil {
		return 0, 0, err
	}

	row = s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM urls")
	var urls int
	if err := row.Scan(&urls); err != nil {
		return 0, 0, err
	}

	return urls, users, nil
}

// generateKey generates a random key with the specified length.
// It uses the characters in the charset string and a random number generator
// to create the key.
//
// The function returns the generated key as a string.
// Example:
//
//	key := generateKey()
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
