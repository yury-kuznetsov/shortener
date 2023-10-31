package main

import (
	"github.com/go-chi/chi"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/app"
	"github.com/yury-kuznetsov/shortener/internal/logger"
	"github.com/yury-kuznetsov/shortener/internal/storage"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
	"net/http"
)

func main() {
	coder := uricoder.NewCoder(storage.NewStorage())
	r := buildRouter(coder)
	config.Init()

	if err := http.ListenAndServe(config.Options.HostAddr, r); err != nil {
		panic(err)
	}
}

func buildRouter(coder *uricoder.Coder) *chi.Mux {
	sugar := logger.NewLogger()

	r := chi.NewRouter()
	r.Get("/{code}", sugar.Handle(handlers.DecodeHandler(coder)))
	r.Post("/api/shorten", sugar.Handle(handlers.EncodeJSONHandler(coder)))
	r.Post("/", sugar.Handle(handlers.EncodeHandler(coder)))
	r.MethodNotAllowed(sugar.Handle(handlers.NotAllowedHandler()))
	//r.Use(middleware.Logger)

	return r
}
