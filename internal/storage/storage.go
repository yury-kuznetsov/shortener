package storage

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type Storage map[string]string

var filename string

func (s *Storage) Get(code string) string {
	return (*s)[code]
}

func (s *Storage) Set(value string) string {
	key := generateKey()
	(*s)[key] = value

	if err := saveToFile(s); err != nil {
		panic(err)
	}
	return key
}

func NewStorage(fName string) *Storage {
	s := make(Storage)
	filename = fName
	if err := loadFromFile(&s); err != nil {
		panic(err)
	}
	return &s
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

func loadFromFile(s *Storage) error {
	if filename == "" {
		return nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil
	}

	if b := len(data); b > 0 {
		err = json.Unmarshal(data, &s)
	}

	return err
}

func saveToFile(s *Storage) error {
	if filename == "" {
		return nil
	}

	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0666)
}
