package uricoder

import (
	"errors"
	"github.com/yury-kuznetsov/shortener/internal/storage"
	"net/url"
)

func NewCoder(s *storage.Storage) *Coder {
	return &Coder{storage: s}
}

type Coder struct {
	storage *storage.Storage
}

func (coder *Coder) ToURI(code string) (string, error) {
	uri := coder.storage.Get(code)
	if uri == "" {
		return "", errors.New("URI not found")
	}
	return uri, nil
}

func (coder *Coder) ToCode(uri string) (string, error) {
	_, err := url.ParseRequestURI(uri)
	if err != nil {
		return "", errors.New("incorrect URI")
	}
	return coder.storage.Set(uri)
}
