package dto

type Response[T any] struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
	Data      T      `json:"data,omitempty"`
}

type PageMeta struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Total  int64 `json:"total"`
}

type PageResponse[T any] struct {
	Meta  PageMeta `json:"meta"`
	Items []T      `json:"items"`
}

type ErrorBody struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func OK[T any](requestID string, data T) Response[T] {
	return Response[T]{
		Code:      0,
		Message:   "ok",
		RequestID: requestID,
		Data:      data,
	}
}

func Err(requestID string, code int, message string) ErrorBody {
	return ErrorBody{
		Code:      code,
		Message:   message,
		RequestID: requestID,
	}
}
