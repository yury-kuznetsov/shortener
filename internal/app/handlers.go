package handlers

import (
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"io"
	"net/http"
	"strings"

	"github.com/yury-kuznetsov/shortener/cmd/storage"
	"github.com/yury-kuznetsov/shortener/cmd/uricoder"
)

func HandlerGet(res http.ResponseWriter, req *http.Request) {
	code := strings.TrimLeft(req.URL.Path, "/")
	uri, err := getCoder().ToURI(code)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	http.Redirect(res, req, uri, http.StatusTemporaryRedirect)
}

func HandlerPost(res http.ResponseWriter, req *http.Request) {
	uri, _ := io.ReadAll(req.Body)
	code, err := getCoder().ToCode(string(uri))
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	res.Header().Set("content-type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	_, _ = res.Write([]byte(config.Options.BaseAddr + "/" + code))
}

func getCoder() *uricoder.Coder {
	return uricoder.NewCoder(storage.ArrStorage)
}
