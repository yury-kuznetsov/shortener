package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yury-kuznetsov/shortener/cmd/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerGet(t *testing.T) {
	arrStorage := storage.ArrStorage

	codes := [2]string{
		arrStorage.Set("https://google.com"),
		arrStorage.Set(""),
	}

	tests := []struct {
		name   string
		code   string
		status int
	}{
		{
			name:   "google",
			code:   codes[0],
			status: http.StatusTemporaryRedirect,
		},
		{
			name:   "bad request",
			code:   codes[1],
			status: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/"+test.code, nil)
			HandlerGet(rec, req)
			res := rec.Result()
			assert.Equal(t, test.status, res.StatusCode)
			defer res.Body.Close()
		})
	}
}

func TestHandlerPost(t *testing.T) {
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
			HandlerPost(rec, req)
			res := rec.Result()
			require.Equal(t, test.status, res.StatusCode)
			defer res.Body.Close()
		})
	}
}
