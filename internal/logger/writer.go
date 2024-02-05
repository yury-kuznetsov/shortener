package logger

import "net/http"

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// Write writes the given byte slice to the underlying ResponseWriter
// and updates the size in the responseData struct.
// It returns the number of bytes written and any error encountered.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader writes the HTTP status code to the underlying ResponseWriter
// and sets the status code in the responseData struct.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
