package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/yury-kuznetsov/shortener/cmd/config"
	"github.com/yury-kuznetsov/shortener/internal/models"
	"github.com/yury-kuznetsov/shortener/internal/uricoder"
)

// DecodeHandler decodes the given code to a URI and redirects the user to the decoded URI.
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

// EncodeBatchHandler handles batch encoding of URLs.
// It receives a JSON array of EncodeBatchRequest and returns a JSON array of EncodeBatchResponse.
// Each EncodeBatchRequest contains an original URL to be encoded.
// Each EncodeBatchResponse contains the correlation ID and the short URL.
// If the encoding fails for any request, EncodeBatchHandler returns an error response.
// The handler function decodes the request body and prepares the response.
// For each request, it calls the ToCode method of the Coder instance to encode the original URL.
// It appends the EncodeBatchResponse to the response array.
// Finally, it sets the content-type header, writes the response array as JSON, and returns the status code 201.
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

// EncodeJSONHandler encodes the given URL to a code and returns the code in a JSON response.
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

// EncodeHandler encodes the given URI and returns the generated code.
// If the code generation fails, it returns a Bad Request error with the corresponding error message.
// It also sets the "content-type" header to "text/plain" and the response status code to StatusCreated if no error occurs.
// Otherwise, it sets the response status code to StatusConflict.
// Finally, it writes the base address concatenated with the generated code to the response body.
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

// PingHandler handles the ping request to check the health of the coder.
// It calls the HealthCheck function of the Coder to perform the check.
// If the health check fails, it returns HTTP StatusInternalServerError.
// Otherwise, it returns HTTP StatusOK.
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

// UserUrlsHandler retrieves and returns a user's URL history in JSON format.
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

// DeleteUrlsHandler deletes URLs based on the given codes.
// It receives a URICoder instance and returns an http.HandlerFunc.
// The handler decodes the incoming JSON request body to get the codes to be deleted.
// It then calls the DeleteUrls method of the URICoder to delete the URLs.
// If there is an error decoding the JSON, it returns a 500 Internal Server Error.
// After deleting the URLs, it sets the response status code to 202 Accepted.
func DeleteUrlsHandler(coder *uricoder.Coder) http.HandlerFunc {
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

// GetStatsHandler retrieves statistics from the `Coder` storage and returns them as a JSON response.
// The returned response includes the number of stored URLs and the number of unique users.
// If an error occurs during the retrieval or encoding of the statistics, an internal server error is returned.
func GetStatsHandler(coder *uricoder.Coder) http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("content-type", "application/json")

		// получаем данные
		urls, users, err := coder.GetStats(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// возвращаем ответ
		response := models.GetStatsResponse{Urls: urls, Users: users}
		if err := json.NewEncoder(res).Encode(response); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusOK)
	}

	return handlerFunc
}

// NotAllowedHandler handles requests that are not allowed.
// If the request method is not GET or POST, it returns a "only GET/POST requests are allowed" error with status code 400.
func NotAllowedHandler() http.HandlerFunc {
	handlerFunc := func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet && req.Method != http.MethodPost {
			http.Error(res, "only GET/POST requests are allowed", http.StatusBadRequest)
		}
	}

	return handlerFunc
}
