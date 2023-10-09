package main

import (
	"net/http"

	"github.com/yury-kuznetsov/shortener/internal/app"
)

func handler(res http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		handlers.HandlerPost(res, req)
		return
	}

	if req.Method == http.MethodGet {
		handlers.HandlerGet(res, req)
		return
	}

	http.Error(res, "only GET/POST requests are allowed", http.StatusBadRequest)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, handler)

	if err := http.ListenAndServe(`:8080`, mux); err != nil {
		panic(err)
	}
}
