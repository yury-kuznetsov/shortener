package models

import "errors"

// EncodeRequest is a struct representing the request for the EncodeJSONHandler method.
// It contains the URL, which represents the original URL to be encoded.
type EncodeRequest struct {
	URL string `json:"url"`
}

// EncodeResponse is a struct representing the response for the EncodeJSONHandler method.
// It contains the Result, which represents the encoded URL.
//
// Usage example:
//
//	response := models.EncodeResponse{Result: config.Options.BaseAddr + "/" + code}
//	if err := json.NewEncoder(res).Encode(response); err != nil {
//	    http.Error(res, err.Error(), http.StatusInternalServerError)
//	    return
//	}
type EncodeResponse struct {
	Result string `json:"result"`
}

// EncodeBatchRequest is a struct representing the request for the EncodeBatchHandler method.
// It contains the CorrelationID, which represents the correlation ID for batch encoding request,
// and the OriginalURL, which represents the original URL to be encoded.
type EncodeBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// EncodeBatchResponse is a struct representing the response for the EncodeBatchHandler method.
// It contains the CorrelationID, which represents the correlation ID for batch encoding request,
// and the ShortURL, which represents the shortened URL.
type EncodeBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// GetByUserResponse is a struct representing the response for the GetByUser method.
// It contains the ShortURL, which represents the shortened URL, and the OriginalURL, which is the original URL.
type GetByUserResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// GetStatsResponse is a struct representing the response for the GetStatsHandler method.
// It contains the number of URLs and users.
type GetStatsResponse struct {
	Urls  int `json:"urls"`
	Users int `json:"users"`
}

// RmvUrlsMsg is a struct representing a message for removing URLs.
// It contains userID, which represents the user ID, and code, which is the code for the URL.
type RmvUrlsMsg struct {
	UserID int
	Code   string
}

// ErrRowDeleted is a variable that represents the error when a row is already deleted.
var ErrRowDeleted = errors.New("запись уже удалена")
