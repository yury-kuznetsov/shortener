package storage

import (
	"math/rand"
	"time"
)

type Storage map[string]string

func (s Storage) Get(code string) string {
	return s[code]
}

func (s Storage) Set(value string) string {
	key := generateKey()
	s[key] = value
	return key
}

var ArrStorage = make(Storage)

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
