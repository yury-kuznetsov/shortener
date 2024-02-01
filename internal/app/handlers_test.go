package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yury-kuznetsov/shortener/internal/storage/memory"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
)

func TestDecodeHandler(t *testing.T) {
	mapStorage := memory.NewStorage()
	coder := uricoder.NewCoder(mapStorage)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	code1, _ := mapStorage.Set(ctx, "https://google.com", 0)
	code2, _ := mapStorage.Set(ctx, "", 0)

	tests := []struct {
		name   string
		code   string
		status int
	}{
		{
			name:   "google",
			code:   code1,
			status: http.StatusTemporaryRedirect,
		},
		{
			name:   "bad request",
			code:   code2,
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/"+test.code, nil)
			DecodeHandler(coder)(rec, req)
			res := rec.Result()
			assert.Equal(t, test.status, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestEncodeHandler(t *testing.T) {
	mapStorage := memory.NewStorage()

	coder := uricoder.NewCoder(mapStorage)
	tests := []struct {
		name   string
		uri    string
		status int
		code   string
	}{
		{
			name:   "google",
			uri:    "https://google.com",
			status: http.StatusCreated,
		},
		{
			name:   "bad request",
			uri:    "incorrect",
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.uri))
			EncodeHandler(coder)(rec, req)
			res := rec.Result()
			require.Equal(t, test.status, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestNotAllowedHandler(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		allowed bool
	}{
		{
			name:    "GET",
			method:  http.MethodGet,
			allowed: true,
		},
		{
			name:    "HEAD",
			method:  http.MethodHead,
			allowed: false,
		},
		{
			name:    "POST",
			method:  http.MethodPost,
			allowed: true,
		},
		{
			name:    "PUT",
			method:  http.MethodPut,
			allowed: false,
		},
		{
			name:    "PATCH",
			method:  http.MethodPatch,
			allowed: false,
		},
		{
			name:    "DELETE",
			method:  http.MethodDelete,
			allowed: false,
		},
		{
			name:    "CONNECT",
			method:  http.MethodConnect,
			allowed: false,
		},
		{
			name:    "OPTIONS",
			method:  http.MethodOptions,
			allowed: false,
		},
		{
			name:    "TRACE",
			method:  http.MethodTrace,
			allowed: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(test.method, "/", nil)
			NotAllowedHandler()(rec, req)
			res := rec.Result()
			require.Equal(t, test.allowed, res.StatusCode != http.StatusBadRequest)
			defer res.Body.Close()
		})
	}
}
