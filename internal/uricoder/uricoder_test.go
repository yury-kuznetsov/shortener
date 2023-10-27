package uricoder

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/yury-kuznetsov/shortener/internal/storage"
	"testing"
)

func TestToURI(t *testing.T) {
	coder := NewCoder(storage.NewStorage())

	codes := [3]string{
		coder.storage.Set("https://google.com"),
		coder.storage.Set("https://ya.ru"),
		coder.storage.Set(""),
	}

	tests := []struct {
		name string
		code string
		uri  string
		err  error
	}{
		{
			name: "google",
			code: codes[0],
			uri:  "https://google.com",
			err:  nil,
		},
		{
			name: "yandex",
			code: codes[1],
			uri:  "https://ya.ru",
			err:  nil,
		},
		{
			name: "not found",
			code: codes[2],
			uri:  "",
			err:  errors.New("URI not found"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri, err := coder.ToURI(test.code)
			assert.Equal(t, uri, test.uri)
			assert.Equal(t, err, test.err)
		})
	}
}

func TestToCode(t *testing.T) {
	coder := NewCoder(storage.NewStorage())

	tests := []struct {
		name string
		uri  string
		code string
		err  error
	}{
		{
			name: "google",
			uri:  "https://google.com/",
			err:  nil,
		},
		{
			name: "yandex",
			uri:  "https://ya.ru",
			err:  nil,
		},
		{
			name: "incorrect",
			uri:  "",
			err:  errors.New("incorrect URI"),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, err := coder.ToCode(test.uri)
			if code != "" {
				uri, _ := coder.ToURI(code)
				assert.Equal(t, uri, test.uri)
			}
			assert.Equal(t, err, test.err)
		})
	}
}
