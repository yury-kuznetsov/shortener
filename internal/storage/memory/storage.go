package memory

import (
	"context"
	"errors"
	"math/rand"
	"time"
)

type Storage map[string]string

func (s *Storage) Get(ctx context.Context, code string) (string, error) {
	v, ok := (*s)[code]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (s *Storage) Set(ctx context.Context, value string) (string, error) {
	key := generateKey()
	(*s)[key] = value
	return key, nil
}

func (s *Storage) HealthCheck(ctx context.Context) error {
	// тут можно добавить проверку занимаемой памяти
	return nil
}

func NewStorage() *Storage {
	return &Storage{}
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
