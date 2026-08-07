package httpx

import "net/http"

type Error struct {
	Ok      bool   `json:"ok"`
	Status  int    `json:"status"`
	Message string `json:"message"`
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

func NewInternalServerError(message string) *Error {
	return NewError(http.StatusInternalServerError, message)
}

func (e *Error) GetStatus() int {
	return e.Status
}
