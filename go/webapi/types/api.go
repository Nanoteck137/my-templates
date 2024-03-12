package types

import (
	"net/http"
)

var (
	ErrNoTodo = NewApiError(http.StatusNotFound, "Todo not found")
	ErrNoChanges = NewApiError(http.StatusBadRequest, "No changes")
)

const (
	StatusSuccess = "success"
	StatusError   = "error"
)

type ApiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  any    `json:"errors,omitempty"`
}

func (err *ApiError) Error() string {
	return err.Message
}

func NewApiError(code int, message string, errors ...any) *ApiError {
	var e any
	if len(errors) > 0 {
		e = errors[0]
	}

	return &ApiError{
		Code:    code,
		Message: message,
		Errors:  e,
	}
}

type ApiResponse struct {
	Status string    `json:"status"`
	Data   any       `json:"data,omitempty"`
	Error  *ApiError `json:"error,omitempty"`
}

func NewApiSuccessResponse(data any) ApiResponse {
	return ApiResponse{
		Status: StatusSuccess,
		Data:   data,
	}
}
