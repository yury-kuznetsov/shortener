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

// Storage represents a map[string]string type, which is used for key-value storage.
// Get retrieves the value associated with the given code and userID from the storage.
// It returns the value and an error if the code is not found in the storage.
// Example usage:
//
//	value, err := storage.Get(ctx, code, userID)
type Storage map[string]string

var filename string

// Get retrieves the value associated with the given code from the Storage instance.
// It takes a context as an argument, which represents the execution context.
// The method expects the code of the value to be retrieved and the userID of the user associated with the value.
// It checks if the code exists in the Storage instance using the code as the key.
// If the code does not exist, it returns an error with the message "not found".
// Otherwise, it returns the value and nil.
//
// Example usage:
// value, err := storage.Get(ctx, code, userID)
//
//	if err != nil {
//	    // handle error
//	}
//
// // use value
func (s *Storage) Get(ctx context.Context, code string, userID int) (string, error) {
	v, ok := (*s)[code]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

// Set adds a new key-value pair to the Storage instance.
// It takes a context as an argument, which represents the execution context.
// The method expects the value to be stored and the userID of the user associated with the value.
// It generates a unique key using the generateKey function and assigns the value to the Storage instance with the generated key.
// After setting the value, it calls the saveToFile function to save the updated Storage instance to a file.
// If an error occurs during the file saving process, it returns an empty string and the error.
// Otherwise, it returns the generated key and nil.
// Example usage:
//
//	key, err := storage.Set(ctx, value, userID)
//	if err != nil {
//	    // handle error
//	}
//	// use key
func (s *Storage) Set(ctx context.Context, value string, userID int) (string, error) {
	key := generateKey()
	(*s)[key] = value
	if err := saveToFile(s); err != nil {
		return "", err
	}
	return key, nil
}

// GetByUser retrieves the data associated with a specific user from the Storage instance.
// It takes a context as an argument, which represents the execution context.
// The method expects the userID of the user whose data needs to be fetched.
// It returns a slice of GetByUserResponse objects, which contain the short URL and original URL.
// If the user has no associated data or if an error occurs, it returns an empty slice and an error respectively.
// Example usage:
//
//	data, err := storage.GetByUser(ctx, userID)
//	if err != nil {
//	    // handle error
//	}
//	// process data
func (s *Storage) GetByUser(ctx context.Context, userID int) ([]models.GetByUserResponse, error) {
	// предстоит переделать хранилище на map[userID][code]value
	return nil, nil
}

// SoftDelete performs a soft delete operation on the Storage instance.
// It takes a context as an argument, which represents the execution context.
// The method expects a slice of messages of type models.RmvUrlsMsg, which contains UserID and Code.
// This method will be modified to change the underlying storage structure to a map[userID][code]value.
// It returns nil if the soft delete operation succeeds, otherwise it returns an error.
// Example usage:
//
//	err := storage.SoftDelete(ctx, messages)
//	if err != nil {
//	    // handle error
//	}
//	// soft delete operation succeeded
func (s *Storage) SoftDelete(ctx context.Context, messages []models.RmvUrlsMsg) error {
	// предстоит переделать хранилище на map[userID][code]value
	return nil
}

// HealthCheck performs a health check on the Storage instance.
// It takes a context as an argument, which represents the execution context.
// This method can be extended to include file presence check and other health checks.
// It returns nil if the health check passes, otherwise it returns an error.
// Example usage:
//
//	err := storage.HealthCheck(ctx)
//	if err != nil {
//	    // handle error
//	}
//	// health check passed
func (s *Storage) HealthCheck(ctx context.Context) error {
	// тут можно добавить проверку наличия файла
	return nil
}

// GetStats retrieves the current statistics of the storage.
func (s *Storage) GetStats(context.Context) (int, int, error) {
	return len(*s), 0, nil
}

// NewStorage creates a new instance of Storage initialized with data from a file.
// It takes a filename as an argument, which represents the file from which the data will be loaded.
// If an error occurs during the loading process, NewStorage returns the error.
// Otherwise, it returns a pointer to the Storage instance and a nil error.
// Example usage:
//
//	storage, err := NewStorage("filename.txt")
//	if err != nil {
//	    // handle error
//	}
//	// use storage instance
//
// Declaration:
//
//	func NewStorage(fName string) (*Storage, error)
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
