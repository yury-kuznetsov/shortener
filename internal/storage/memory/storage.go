package memory

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/yury-kuznetsov/shortener/internal/models"
)

// Storage represents a map-based storage that stores key-value pairs. The key is a string, and the value is also a string.
type Storage map[string]string

// Get retrieves the value associated with the given code from the storage.
// If the code is not found in the storage, it returns an empty string and an error message "not found".
// Example usage:
//
//	storage := NewStorage()
//	ctx := context.Background()
//	value, err := storage.Get(ctx, "code123", userID)
//	if err != nil {
//	    // handle error
//	}
//	// use value
func (s *Storage) Get(ctx context.Context, code string, userID int) (string, error) {
	v, ok := (*s)[code]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

// Set adds a new key-value pair to the storage.
// The key is generated using the generateKey function.
// The value is stored in the storage using the generated key.
// The method returns the generated key and nil error.
// Example usage:
//
//	storage := NewStorage()
//	ctx := context.Background()
//	key, err := storage.Set(ctx, "https://site.com", userID)
//	if err != nil {
//	    // handle error
//	}
//	// use key
func (s *Storage) Set(ctx context.Context, value string, userID int) (string, error) {
	key := generateKey()
	(*s)[key] = value
	return key, nil
}

// GetByUser retrieves the data for a specific user from the storage.
// Currently, the storage is implemented as a map[string]string, but it will be refactored to a map[userID][code]value structure in the future.
// The method returns an array of models.GetByUserResponse and an error.
// Example usage:
//
//	storage := NewStorage()
//	ctx := context.Background()
//	data, err := storage.GetByUser(ctx, userID)
//	if err != nil {
//		// handle error
//	}
//	// use data
//
//	type models.GetByUserResponse struct {
//		ShortURL    string `json:"short_url"`
//		OriginalURL string `json:"original_url"`
//	}
func (s *Storage) GetByUser(ctx context.Context, userID int) ([]models.GetByUserResponse, error) {
	// предстоит переделать хранилище на map[userID][code]value
	return nil, nil
}

// SoftDelete deletes the specified messages from the storage.
// Currently, the storage is implemented as a map[string]string, but it will be refactored to a map[userID][code]value structure in the future.
func (s *Storage) SoftDelete(ctx context.Context, messages []models.RmvUrlsMsg) error {
	// предстоит переделать хранилище на map[userID][code]value
	return nil
}

// HealthCheck checks the health of the storage.
// It can include memory utilization checks.
func (s *Storage) HealthCheck(ctx context.Context) error {
	// тут можно добавить проверку занимаемой памяти
	return nil
}

// NewStorage creates a new instance of the Storage struct.
func NewStorage() *Storage {
	return &Storage{}
}

// GetStats retrieves the current statistics of the storage.
func (s *Storage) GetStats(context.Context) (int, int, error) {
	return len(*s), 0, nil
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
