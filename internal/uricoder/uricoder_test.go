package uricoder

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yury-kuznetsov/shortener/internal/storage/file"
	"github.com/yury-kuznetsov/shortener/internal/storage/memory"
)

func TestToURI(t *testing.T) {
	s, err := file.NewStorage("")
	require.NoError(t, err)

	coder := NewCoder(s)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	code1, _ := s.Set(ctx, "https://google.com", 0)
	code2, _ := s.Set(ctx, "https://ya.ru", 0)
	code3, _ := s.Set(ctx, "", 0)

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
			uri, err := coder.ToURI(ctx, test.code, 0)
			assert.Equal(t, uri, test.uri)
			assert.Equal(t, err, test.err)
		})
	}
}

func TestToCode(t *testing.T) {
	s, err := file.NewStorage("")
	require.NoError(t, err)

	coder := NewCoder(s)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
			code, err := coder.ToCode(ctx, test.uri, 0)
			if code != "" {
				uri, _ := coder.ToURI(ctx, code, 0)
				assert.Equal(t, uri, test.uri)
			}
			assert.Equal(t, err, test.err)
		})
	}
}

func BenchmarkToURI(b *testing.B) {
	s := memory.NewStorage()
	code, _ := s.Set(context.Background(), "https://ya.ru", 0)
	coder := NewCoder(s)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = coder.ToURI(context.Background(), code, 0)
	}
}

func BenchmarkToCode(b *testing.B) {
	s := memory.NewStorage()
	coder := NewCoder(s)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = coder.ToCode(context.Background(), "https://ya.ru", 0)
	}
}
