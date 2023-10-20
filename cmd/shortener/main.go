package main

import (
	"github.com/go-chi/chi"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/app"
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
	r := chi.NewRouter()
	r.Get("/{code}", handlers.DecodeHandler(coder))
	r.Post("/", handlers.EncodeHandler(coder))
	r.MethodNotAllowed(handlers.NotAllowedHandler())

	return r
}
