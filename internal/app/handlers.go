package handlers

import (
	"encoding/json"
	"errors"
	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/models"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func DecodeHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		userID, err := strconv.Atoi(req.Header.Get("Content-User-ID"))
		if err != nil {
			userID = 0
		}

		code := strings.TrimLeft(req.URL.Path, "/")
		uri, err := coder.ToURI(req.Context(), code, userID)
		if err != nil {
			if errors.Is(err, models.ErrRowDeleted) {
				res.WriteHeader(http.StatusGone)
				return
			}
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
		userID, err := strconv.Atoi(req.Header.Get("Content-User-ID"))
		if err != nil {
			userID = 0
		}

		// готовим ответ
		var response []models.EncodeBatchResponse
		for _, v := range request {
			code, err := coder.ToCode(req.Context(), v.OriginalURL, userID)
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
		userID, err := strconv.Atoi(req.Header.Get("Content-User-ID"))
		if err != nil {
			userID = 0
		}

		// запускаем обработку
		code, err := coder.ToCode(req.Context(), request.URL, userID)
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
		userID, err := strconv.Atoi(req.Header.Get("Content-User-ID"))
		if err != nil {
			userID = 0
		}
		uri, _ := io.ReadAll(req.Body)
		code, err := coder.ToCode(req.Context(), string(uri), userID)
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
		err := coder.HealthCheck(req.Context())
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}

	return handlerFunc
}

func UserUrlsHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", "application/json")

		// проверяем авторизацию
		userID, err := strconv.Atoi(req.Header.Get("Content-User-ID"))
		if err != nil {
			userID = 0
		}
		if userID == 0 {
			res.WriteHeader(http.StatusUnauthorized)
			return
		}

		// запускаем обработку запроса
		data, err := coder.GetHistory(req.Context(), userID)
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(data) == 0 {
			res.WriteHeader(http.StatusNoContent)
			return
		}

		// возвращаем ответ
		var response []models.GetByUserResponse
		for _, v := range data {
			response = append(response, models.GetByUserResponse{
				ShortURL:    config.Options.BaseAddr + "/" + v.ShortURL,
				OriginalURL: v.OriginalURL,
			})
		}
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}

	return handlerFunc
}

func DeleteUrls(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		userID, err := strconv.Atoi(req.Header.Get("Content-User-ID"))
		if err != nil {
			userID = 0
		}

		var codes []string
		if err := json.NewDecoder(req.Body).Decode(&codes); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		_ = coder.DeleteUrls(codes, userID)

		res.WriteHeader(http.StatusAccepted)
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
