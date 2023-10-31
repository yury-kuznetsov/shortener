package models

type EncodeRequest struct {
	URL string `json:"url"`
}

type EncodeResponse struct {
	Result string `json:"result"`
}
