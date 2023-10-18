package main

import (
	"errors"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yury-kuznetsov/shortener/internal/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type request struct {
	method string
	target string
	body   string
}
type response struct {
	status int
	error  error
}
type testCase struct {
	name     string
	request  request
	response response
}

func TestRequests(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/{code}", ToURI)
	r.Post("/", ToCode)
	r.MethodNotAllowed(NotAllowed)

	ts := httptest.NewServer(r)
	defer ts.Close()

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	tests := []testCase{
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
			testRequest(t, ts, test)
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, test testCase) {
	req, err := http.NewRequest(
		test.request.method,
		ts.URL+test.request.target,
		strings.NewReader(test.request.body),
	)
	require.NoError(t, err)

	resp, err := ts.Client().Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	assert.Equal(t, test.response.status, resp.StatusCode)

	if test.response.error != nil {
		assert.Equal(t, test.response.error.Error(), string(respBody))
	}
}
