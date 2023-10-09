package main

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/yury-kuznetsov/shortener/cmd/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler(t *testing.T) {
	type request struct {
		method string
		target string
		body   string
	}
	type response struct {
		status int
		error  error
	}
	tests := []struct {
		name     string
		request  request
		response response
	}{
		{
			name: "GET-request",
			request: request{
				method: http.MethodGet,
				target: "/" + storage.ArrStorage.Set("https://google.com"),
				body:   "",
			},
			response: response{
				status: http.StatusTemporaryRedirect,
				error:  nil,
			},
		},
		{
			name: "POST-request",
			request: request{
				method: http.MethodPost,
				target: "/",
				body:   "https://google.com",
			},
			response: response{
				status: http.StatusCreated,
				error:  nil,
			},
		},
		{
			name: "PUT-request",
			request: request{
				method: http.MethodPut,
				target: "/",
				body:   "",
			},
			response: response{
				status: http.StatusBadRequest,
				error:  errors.New("only GET/POST requests are allowed\n"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(test.request.method, test.request.target, strings.NewReader(test.request.body))
			handler(rec, req)
			res := rec.Result()
			assert.Equal(t, test.response.status, res.StatusCode)
			if test.response.error != nil {
				assert.Equal(t, test.response.error.Error(), rec.Body.String())
			}
			defer res.Body.Close()
		})
	}
}
