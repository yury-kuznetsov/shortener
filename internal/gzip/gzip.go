package gzip

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

// Header returns the header map of the underlying http.ResponseWriter (c.w).
// The header map represents the key-value pairs of the HTTP response header.
// Any changes made to the returned map will be reflected in the response.
// Usage example:
//
//	cw := &compressWriter{w: responseWriter}
//	header := cw.Header()
//	header.Set("Content-Type", "application/json")
//	header.Add("X-Custom-Header", "value")
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Write writes the given byte slice to the underlying gzip.Writer (c.zw).
// It returns the number of bytes written and any error encountered during the write operation.
func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

// WriteHeader sets the HTTP response header for the given status code.
// If the status code is less than 300, it also sets the "Content-Encoding" header to "gzip".
// The method then calls the WriteHeader method of the underlying http.ResponseWriter (c.w).
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Close closes the gzip writer (c.zw) and returns any error encountered during the close operation.
func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

// Read reads data from the gzip reader (c.zr) and stores it in the provided byte slice (p).
// It returns the number of bytes read (n) and any error encountered during the read (err).
// The provided byte slice (p) must have sufficient capacity to hold the data read.
func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

// Close closes the compressReader by first closing the underlying reader (c.r) and then closing the gzip reader (c.zr).
// It returns any error encountered while closing either reader.
func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

// Handle wraps an HTTP handler function and adds functionality for handling compression.
// It checks if the client supports compressed data and if the client sends compressed data,
// then it sets up the appropriate writer or reader to handle compression.
// It then calls the original handler function with the modified writer and reader.
func Handle(handler http.HandlerFunc) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		ow := res

		// проверяем, что клиент умеет получать сжатые данные
		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(res)
			ow = cw
			defer cw.Close()
		}

		// проверяем, что клиент отправил сжатые данные
		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				fmt.Println(err.Error())
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = cr
			defer cr.Close()
		}

		handler(ow, req)
	}

	return handlerFunc
}
