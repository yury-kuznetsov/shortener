// Package logger provides a simple logging functionality.
//
// The logger package allows logging messages with different log levels.
//
// All log messages are written to the standard output.
//
// Example usage:
//
//	logger.Info("This is an informational message")
//	logger.Debug("This is a debug message")
//	logger.Error("This is an error message")
package logger

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// Logger is a type that wraps around zap.SugaredLogger for logging purposes.
type Logger struct {
	sugar zap.SugaredLogger
}

// NewLogger returns a new instance of the Logger struct.
func NewLogger() *Logger {
	loggerZap, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize logger")
	}
	defer loggerZap.Sync()

	logger := Logger{sugar: *loggerZap.Sugar()}

	return &logger
}

// Handle method handles the HTTP request and response.
func (l *Logger) Handle(handler http.HandlerFunc) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()
		requestURI := req.RequestURI
		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: res,
			responseData:   responseData,
		}

		handler(&lw, req)

		duration := time.Since(start)
		l.sugar.Infoln(
			"uri", requestURI,
			"method", req.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}

	return handlerFunc
}
