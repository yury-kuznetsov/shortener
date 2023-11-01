package main

import (
	"github.com/go-chi/chi"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/app"
	"github.com/yury-kuznetsov/shortener/internal/gzip"
	"github.com/yury-kuznetsov/shortener/internal/logger"
	"github.com/yury-kuznetsov/shortener/internal/storage"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
	"net/http"
)

func main() {
	config.Init()
	coder := uricoder.NewCoder(storage.NewStorage(config.Options.FilePath))
	r := buildRouter(coder)

	if err := http.ListenAndServe(config.Options.HostAddr, r); err != nil {
		panic(err)
	}
}

func buildRouter(coder *uricoder.Coder) *chi.Mux {
	sugar := logger.NewLogger()

	r := chi.NewRouter()
	r.Get("/{code}", gzip.Handle(sugar.Handle(handlers.DecodeHandler(coder))))
	r.Post("/api/shorten", gzip.Handle(sugar.Handle(handlers.EncodeJSONHandler(coder))))
	r.Post("/", gzip.Handle(sugar.Handle(handlers.EncodeHandler(coder))))
	r.MethodNotAllowed(gzip.Handle(sugar.Handle(handlers.NotAllowedHandler())))
	//r.Use(middleware.Logger)

	return r
}
