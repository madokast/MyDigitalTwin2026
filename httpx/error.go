package httpx

import (
	"net/http"
)

type Error struct {
	Ok      bool   `json:"ok" jsonschema:"Whether the request succeeded"`
	Status  int    `json:"status" jsonschema:"HTTP status code"`
	Message string `json:"message" jsonschema:"Error message"`
}

func NewError(status int, message string) *Error {
	return &Error{
		Ok:      false,
		Status:  status,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *Error {
	return NewError(http.StatusUnauthorized, message)
}

func NewBadRequestError(message string) *Error {
	return NewError(http.StatusBadRequest, message)
}

func NewNotFoundError(message string) *Error {
	return NewError(http.StatusNotFound, message)
}

func NewInternalServerError(message string) *Error {
	return NewError(http.StatusInternalServerError, message)
}

func (e Error) GetStatus() int {
	return e.Status
}
