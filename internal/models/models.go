package models

type EncodeRequest struct {
	URL string `json:"url"`
}

type EncodeResponse struct {
	Result string `json:"result"`
}

type EncodeBatchRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type EncodeBatchResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type GetByUserResponse struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type RmvUrlsMsg struct {
	UserID int
	Code   string
}
