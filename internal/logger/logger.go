package logger

import (
	"go.uber.org/zap"
	"net/http"
	"time"
)

type Logger struct {
	sugar zap.SugaredLogger
}

func NewLogger() *Logger {
	loggerZap, err := zap.NewDevelopment()
	if err != nil {
		panic("cannot initialize logger")
	}
	defer loggerZap.Sync()

	logger := Logger{sugar: *loggerZap.Sugar()}

	return &logger
}

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
