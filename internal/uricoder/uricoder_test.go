package uricoder

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yury-kuznetsov/shortener/internal/storage/file"
	"testing"
)

func TestToURI(t *testing.T) {
	s, err := file.NewStorage("")
	require.NoError(t, err)

	coder := NewCoder(s)

	code1, _ := s.Set("https://google.com")
	code2, _ := s.Set("https://ya.ru")
	code3, _ := s.Set("")

	tests := []struct {
		name string
		code string
		uri  string
		err  error
	}{
		{
			name: "google",
			code: code1,
			uri:  "https://google.com",
			err:  nil,
		},
		{
			name: "yandex",
			code: code2,
			uri:  "https://ya.ru",
			err:  nil,
		},
		{
			name: "not found",
			code: code3,
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
	s, err := file.NewStorage("")
	require.NoError(t, err)

	coder := NewCoder(s)

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
