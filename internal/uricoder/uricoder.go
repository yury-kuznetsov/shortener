package uricoder

import (
	"context"
	"errors"
	"net/url"
)

func NewCoder(s Storage) *Coder {
	return &Coder{storage: s}
}

type Coder struct {
	storage Storage
}

func (coder *Coder) ToURI(ctx context.Context, code string) (string, error) {
	uri, err := coder.storage.Get(ctx, code)
	if err != nil {
		return "", err
	}
	if uri == "" {
		return "", errors.New("URI not found")
	}
	return uri, nil
}

func (coder *Coder) ToCode(ctx context.Context, uri string) (string, error) {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", errors.New("incorrect URI")
	}
	return coder.storage.Set(ctx, uri)
}

func (coder *Coder) HealthCheck(ctx context.Context) error {
	return coder.storage.HealthCheck(ctx)
}
