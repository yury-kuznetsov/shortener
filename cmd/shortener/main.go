package main

import (
	"github.com/go-chi/chi"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/app"
	"github.com/yury-kuznetsov/shortener/internal/auth"
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
	r.Get("/{code}", auth.Handle(gzip.Handle(sugar.Handle(handlers.DecodeHandler(coder))), true))
	r.Get("/ping", auth.Handle(gzip.Handle(sugar.Handle(handlers.PingHandler(coder))), true))
	r.Get("/api/user/urls", auth.Handle(gzip.Handle(sugar.Handle(handlers.UserUrlsHandler(coder))), false))
	r.Post("/api/shorten/batch", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeBatchHandler(coder))), true))
	r.Post("/api/shorten", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeJSONHandler(coder))), true))
	r.Post("/", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeHandler(coder))), true))
	r.MethodNotAllowed(auth.Handle(gzip.Handle(sugar.Handle(handlers.NotAllowedHandler())), true))

	return r
}
