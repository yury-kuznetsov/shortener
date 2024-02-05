package main

import (
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	handlers "github.com/yury-kuznetsov/shortener/internal/app"
	"github.com/yury-kuznetsov/shortener/internal/auth"
	"github.com/yury-kuznetsov/shortener/internal/gzip"
	"github.com/yury-kuznetsov/shortener/internal/logger"
	"github.com/yury-kuznetsov/shortener/internal/storage/database"
	"github.com/yury-kuznetsov/shortener/internal/storage/file"
	"github.com/yury-kuznetsov/shortener/internal/storage/memory"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
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
	r.Delete("/api/user/urls", auth.Handle(gzip.Handle(sugar.Handle(handlers.DeleteUrlsHandler(coder))), true))
	r.Post("/api/shorten/batch", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeBatchHandler(coder))), true))
	r.Post("/api/shorten", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeJSONHandler(coder))), true))
	r.Post("/", auth.Handle(gzip.Handle(sugar.Handle(handlers.EncodeHandler(coder))), true))
	r.MethodNotAllowed(auth.Handle(gzip.Handle(sugar.Handle(handlers.NotAllowedHandler())), true))

	// обработчики для pprof
	r.Handle("/debug/pprof/*", http.HandlerFunc(pprof.Index))
	r.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	r.Handle("/debug/pprof/heap", pprof.Handler("heap"))

	// обработчик для favicon.ico (иначе перехватит DecodeHandler)
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	return r
}
