package handlers

import (
	"encoding/json"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/models"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
	"io"
	"net/http"
	"strings"
)

func DecodeHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		code := strings.TrimLeft(req.URL.Path, "/")
		uri, err := coder.ToURI(code)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(res, req, uri, http.StatusTemporaryRedirect)
	}

	return handlerFunc
}

func EncodeJSONHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		// принимаем запрос
		var request models.EncodeRequest
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// запускаем обработку
		code, err := coder.ToCode(request.URL)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		// возвращаем ответ
		response := models.EncodeResponse{Result: config.Options.BaseAddr + "/" + code}
		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return handlerFunc
}

func EncodeHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		uri, _ := io.ReadAll(req.Body)
		code, err := coder.ToCode(string(uri))
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
		res.Header().Set("content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		_, _ = res.Write([]byte(config.Options.BaseAddr + "/" + code))
	}

	return handlerFunc
}

func NotAllowedHandler() http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet && req.Method != http.MethodPost {
			http.Error(res, "only GET/POST requests are allowed", http.StatusBadRequest)
		}
	}

	return handlerFunc
}
