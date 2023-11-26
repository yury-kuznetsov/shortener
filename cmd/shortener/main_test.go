package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yury-kuznetsov/shortener/internal/storage/memory"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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
	s := memory.NewStorage()

	coder := uricoder.NewCoder(s)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ts := httptest.NewServer(buildRouter(coder))
	defer ts.Close()

	client := ts.Client()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	code, _ := coder.ToCode(ctx, "https://google.com", 0)

	tests := []testCase{
		{
			name: "GET-request",
			request: request{
				method: http.MethodGet,
				target: "/" + code,
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
			name: "POST-json-request",
			request: request{
				method: http.MethodPost,
				target: "/api/shorten",
				body:   "{\"url\": \"https://ya.ru\"}",
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

	// махинизм сжатие проверим отдельно
	req.Header.Set("Accept-Encoding", "")
	req.Header.Set("Content-Encoding", "")

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

func TestGzipCompression(t *testing.T) {
	s := memory.NewStorage()

	coder := uricoder.NewCoder(s)
	ts := httptest.NewServer(buildRouter(coder))
	defer ts.Close()

	t.Run("sends_gzip", func(t *testing.T) {
		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		_, err := zb.Write([]byte("https://google.com"))
		require.NoError(t, err)
		err = zb.Close()
		require.NoError(t, err)

		r := httptest.NewRequest("POST", ts.URL+"/", buf)
		r.RequestURI = ""
		r.Header.Set("Content-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		_, err = io.ReadAll(resp.Body)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})

	t.Run("accepts_gzip", func(t *testing.T) {
		buf := bytes.NewBufferString("https://ya.ru")
		r := httptest.NewRequest("POST", ts.URL+"/", buf)
		r.RequestURI = ""
		r.Header.Set("Accept-Encoding", "gzip")

		resp, err := http.DefaultClient.Do(r)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)

		defer resp.Body.Close()

		zr, err := gzip.NewReader(resp.Body)
		require.NoError(t, err)

		_, err = io.ReadAll(zr)
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp.StatusCode)
	})
}
