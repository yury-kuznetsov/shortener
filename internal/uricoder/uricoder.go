package uricoder

import (
	"errors"
	"net/url"
)

func NewCoder(s Storage) *Coder {
	return &Coder{storage: s}
}

type Coder struct {
	storage Storage
}

func (coder *Coder) ToURI(code string) (string, error) {
	uri, err := coder.storage.Get(code)
	if err != nil {
		return "", err
	}
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

func (coder *Coder) HealthCheck() error {
	return coder.storage.HealthCheck()
}
