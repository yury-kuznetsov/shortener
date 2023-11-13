package database

import (
	_ "github.com/jackc/pgx/v5/stdlib"
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
