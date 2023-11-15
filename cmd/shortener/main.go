package main

import (
	"github.com/go-chi/chi"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/app"
	"github.com/yury-kuznetsov/shortener/internal/gzip"
	"github.com/yury-kuznetsov/shortener/internal/logger"
	"github.com/yury-kuznetsov/shortener/internal/storage/database"
	"github.com/yury-kuznetsov/shortener/internal/storage/file"
	"github.com/yury-kuznetsov/shortener/internal/storage/memory"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
	"net/http"
)

func main() {
	config.Init()

	storage, err := buildStorage()
	if err != nil {
		panic(err)
	}
	coder := uricoder.NewCoder(storage)

	r := buildRouter(coder)

	if err := http.ListenAndServe(config.Options.HostAddr, r); err != nil {
		panic(err)
	}
}

func buildStorage() (uricoder.Storage, error) {
	if len(config.Options.Database) > 0 {
		return database.NewStorage(config.Options.Database)
	}
	if len(config.Options.FilePath) > 0 {
		return file.NewStorage(config.Options.FilePath)
	}
	return memory.NewStorage(), nil
}

func buildRouter(coder *uricoder.Coder) *chi.Mux {
	sugar := logger.NewLogger()

	r := chi.NewRouter()
	r.Get("/{code}", gzip.Handle(sugar.Handle(handlers.DecodeHandler(coder))))
	r.Get("/ping", gzip.Handle(sugar.Handle(handlers.PingHandler(coder))))
	r.Post("/api/shorten/batch", gzip.Handle(sugar.Handle(handlers.EncodeBatchHandler(coder))))
	r.Post("/api/shorten", gzip.Handle(sugar.Handle(handlers.EncodeJSONHandler(coder))))
	r.Post("/", gzip.Handle(sugar.Handle(handlers.EncodeHandler(coder))))
	r.MethodNotAllowed(gzip.Handle(sugar.Handle(handlers.NotAllowedHandler())))
	//r.Use(middleware.Logger)

	return r
}
