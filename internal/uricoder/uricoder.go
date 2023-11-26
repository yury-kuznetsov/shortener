package uricoder

import (
	"context"
	"errors"
	"github.com/yury-kuznetsov/shortener/internal/models"
	"net/url"
)

func NewCoder(s Storage) *Coder {
	return &Coder{storage: s}
}

type Coder struct {
	storage Storage
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
