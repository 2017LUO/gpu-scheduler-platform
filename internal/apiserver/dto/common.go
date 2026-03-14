package dto

type Response[T any] struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Data      T      `json:"data"`
}

type ListMeta struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Count  int `json:"count"`
}

type ListData[T any] struct {
	Items []T      `json:"items"`
	Meta  ListMeta `json:"meta"`
}

type ErrorData struct {
	Path string `json:"path,omitempty"`
}

type Empty struct{}
