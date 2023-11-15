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

func EncodeBatchHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		// принимаем запрос
		var request []models.EncodeBatchRequest
		if err := json.NewDecoder(req.Body).Decode(&request); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// готовим ответ
		var response []models.EncodeBatchResponse
		for _, v := range request {
			code, err := coder.ToCode(v.OriginalURL)
			if err != nil {
				http.Error(res, err.Error(), http.StatusBadRequest)
				return
			}

			response = append(response, models.EncodeBatchResponse{
				CorrelationID: v.CorrelationID,
				ShortURL:      config.Options.BaseAddr + "/" + code,
			})
		}

		// возвращаем результат
		res.Header().Set("content-type", "application/json")
		res.WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
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
		if code == "" && err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		// возвращаем ответ
		response := models.EncodeResponse{Result: config.Options.BaseAddr + "/" + code}
		res.Header().Set("content-type", "application/json")
		if err != nil {
			res.WriteHeader(http.StatusConflict)
		} else {
			res.WriteHeader(http.StatusCreated)
		}
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	return handlerFunc
}

func EncodeHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		// обрабатываем запрос
		uri, _ := io.ReadAll(req.Body)
		code, err := coder.ToCode(string(uri))
		if code == "" && err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		// возвращаем ответ
		res.Header().Set("content-type", "text/plain")
		if err != nil {
			res.WriteHeader(http.StatusConflict)
		} else {
			res.WriteHeader(http.StatusCreated)
		}
		_, _ = res.Write([]byte(config.Options.BaseAddr + "/" + code))
	}

	return handlerFunc
}

func PingHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		err := coder.HealthCheck()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
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
