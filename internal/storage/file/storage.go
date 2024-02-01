package file

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/yury-kuznetsov/shortener/internal/models"
)

type Storage map[string]string

var filename string

func (s *Storage) Get(ctx context.Context, code string, userID int) (string, error) {
	v, ok := (*s)[code]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (s *Storage) Set(ctx context.Context, value string, userID int) (string, error) {
	key := generateKey()
	(*s)[key] = value
	if err := saveToFile(s); err != nil {
		return "", err
	}
	return key, nil
}

func (s *Storage) GetByUser(ctx context.Context, userID int) ([]models.GetByUserResponse, error) {
	// предстоит переделать хранилище на map[userID][code]value
	return nil, nil
}

func (s *Storage) SoftDelete(ctx context.Context, messages []models.RmvUrlsMsg) error {
	// предстоит переделать хранилище на map[userID][code]value
	return nil
}

func (s *Storage) HealthCheck(ctx context.Context) error {
	// тут можно добавить проверку наличия файла
	return nil
}

func NewStorage(fName string) (*Storage, error) {
	s := make(Storage)
	filename = fName
	if err := loadFromFile(&s); err != nil {
		return &s, err
	}
	return &s, nil
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
