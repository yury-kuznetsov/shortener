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

func NewCoder(s Storage) *Coder {
	instance := &Coder{
		storage:     s,
		rmvUrlsChan: make(chan models.RmvUrlsMsg, 1024),
	}

	go instance.rmvUrls()

	return instance
}

type Coder struct {
	storage     Storage
	rmvUrlsChan chan models.RmvUrlsMsg
}

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

func (coder *Coder) ToCode(ctx context.Context, uri string, userID int) (string, error) {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", errors.New("incorrect URI")
	}
	return coder.storage.Set(ctx, uri, userID)
}

func (coder *Coder) GetHistory(ctx context.Context, userID int) ([]models.GetByUserResponse, error) {
	return coder.storage.GetByUser(ctx, userID)
}

func (coder *Coder) HealthCheck(ctx context.Context) error {
	return coder.storage.HealthCheck(ctx)
}

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
