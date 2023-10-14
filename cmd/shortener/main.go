package main

import (
	"github.com/go-chi/chi"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/app"
	"net/http"
)

func ToURI(w http.ResponseWriter, r *http.Request) {
	handlers.HandlerGet(w, r)
}

func ToCode(w http.ResponseWriter, r *http.Request) {
	handlers.HandlerPost(w, r)
}

func NotAllowed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodPost {
		http.Error(w, "only GET/POST requests are allowed", http.StatusBadRequest)
	}
}

func main() {
	r := chi.NewRouter()
	r.Get("/{code}", ToURI)
	r.Post("/", ToCode)
	r.MethodNotAllowed(NotAllowed)

	config.Init()

	if err := http.ListenAndServe(config.Options.HostAddr, r); err != nil {
		panic(err)
	}
}
