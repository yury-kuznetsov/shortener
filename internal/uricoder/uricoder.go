// Package uricoder is a utility package that provides functions to encode and decode URI strings
package uricoder

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/yury-kuznetsov/shortener/internal/models"
)

// NewCoder initializes a new instance of the Coder struct with the provided Storage implementation.
// It creates a new Coder instance and sets the storage field to the provided Storage implementation.
// It also creates a new channel rmvUrlsChan with a buffer size of 1024 and assigns it to the rmvUrlsChan field.
// It then starts a goroutine to handle the rmvUrls channel.
// The NewCoder function returns the new instance of the Coder struct.
func NewCoder(s Storage) *Coder {
	instance := &Coder{
		storage:     s,
		rmvUrlsChan: make(chan models.RmvUrlsMsg, 1024),
	}

	go instance.rmvUrls()

	return instance
}

// Coder is a struct that represents a coder object. It contains a storage object for data storage and a channel for removing URLs.
//
// Declaration:
//
//	type Coder struct {
//	    storage     Storage
//	    rmvUrlsChan chan models.RmvUrlsMsg
//	}
//
// Usage Example 1:
//
//	c := NewCoder(storageInstance)
//
// Usage Example 2:
//
//	uri, err := c.ToURI(ctx, code, userID)
//
// Usage Example 3:
//
//	code, err := c.ToCode(ctx, uri, userID)
//
// Usage Example 4:
//
//	history, err := c.GetHistory(ctx, userID)
//
// Usage Example 5:
//
//	err := c.HealthCheck(ctx)
//
// Usage Example 6:
//
//	err := c.DeleteUrls(codes, userID)
//
// Usage Example 7:
//
//	c.rmvUrls()
//
// Related Declarations:
// - Storage: an interface for data storage
//
// Related structs:
// - models.RmvUrlsMsg: a struct representing a message for removing URLs
//
// Related methods:
// - ToURI: retrieves the original URL from the provided code and user ID
// - ToCode: retrieves the code for the provided original URL and user ID
// - GetHistory: retrieves the history of URLs for the provided user ID
// - HealthCheck: checks the health of the storage
// - DeleteUrls: deletes multiple URLs for the provided codes and user ID
// - rmvUrls: removes URLs from the storage based on messages received through the rmvUrlsChan channel
type Coder struct {
	storage     Storage
	rmvUrlsChan chan models.RmvUrlsMsg
}

// ToURI returns the URI associated with the given code and user ID.
// It retrieves the URI from the storage using the provided context,
// and returns an error if the code is not found or if there is an error
// retrieving the URI from the storage.
func (coder *Coder) ToURI(ctx context.Context, code string, userID int) (string, error) {
	uri, err := coder.storage.Get(ctx, code, userID)
	if err != nil {
		return "", err
	}
	if uri == "" {
		return "", errors.New("URI not found")
	}
	return uri, nil
}

// ToCode returns the code associated with the given URI and user ID.
// It parses the URI using the url.ParseRequestURI method and returns an error
// if the URI is incorrect. Otherwise, it sets the URI in the storage using
// the provided context and user ID.
func (coder *Coder) ToCode(ctx context.Context, uri string, userID int) (string, error) {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", errors.New("incorrect URI")
	}
	return coder.storage.Set(ctx, uri, userID)
}

// GetHistory returns the history of URLs for a given user.
// It retrieves the URLs from the storage using the provided context and user ID.
// It returns a slice of models.GetByUserResponse, which contains the short URL and original URL.
// If there is an error retrieving the history or if the history is empty, it returns an error.
// Example usage:
//
//	data, err := coder.GetHistory(req.Context(), userID)
//	if err != nil {
//	    // handle error
//	}
//	if len(data) == 0 {
//	    // handle empty history
//	}
//	// process data
func (coder *Coder) GetHistory(ctx context.Context, userID int) ([]models.GetByUserResponse, error) {
	return coder.storage.GetByUser(ctx, userID)
}

// HealthCheck checks the health of the storage.
// It delegates the health check to the storage implementation by calling the HealthCheck method on the storage.
// If there is an error while performing the health check, it is returned.
// Otherwise, it returns nil.
// Usage example:
//
//	err := coder.HealthCheck(req.Context())
//	if err != nil {
//		res.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//	res.WriteHeader(http.StatusOK)
func (coder *Coder) HealthCheck(ctx context.Context) error {
	return coder.storage.HealthCheck(ctx)
}

// DeleteUrls deletes multiple URLs associated with the given codes and user ID.
// It sends a message to the rmvUrlsChan for each code to be deleted,
// triggering the actual deletion process in the background.
// The userID is printed to the console for debugging purposes.
// This method returns nil as there is no error handling in this implementation.
func (coder *Coder) DeleteUrls(codes []string, userID int) error {
	fmt.Println("userID: " + strconv.Itoa(userID))
	for _, code := range codes {
		coder.rmvUrlsChan <- models.RmvUrlsMsg{UserID: userID, Code: code}
	}

	return nil
}

func (coder *Coder) rmvUrls() {
	ticker := time.NewTicker(10 * time.Second)

	var messages []models.RmvUrlsMsg

	for {
		select {
		case message := <-coder.rmvUrlsChan:
			messages = append(messages, message)
		case <-ticker.C:
			if len(messages) == 0 {
				continue
			}
			err := coder.storage.SoftDelete(context.TODO(), messages)
			if err != nil {
				fmt.Print(err)
				continue
			}
			messages = nil
		}
	}
}
