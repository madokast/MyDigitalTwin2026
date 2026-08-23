package httpx

import (
	"net/http"
	"sync"

	"github.com/google/jsonschema-go/jsonschema"
)

type Error struct {
	Ok      FalseType `json:"ok" jsonschema:"Whether the request succeeded"`
	Status  int    `json:"status" jsonschema:"HTTP status code"`
	Message string `json:"message" jsonschema:"Error message"`
}

func NewError(status int, message string) *Error {
	return &Error{
		Ok:      False,
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

func ErrorSchema() *jsonschema.Schema {
	return errorSchema()
}

var errorSchema = sync.OnceValue(func() *jsonschema.Schema {
	s, err := jsonschema.For[Error](&jsonschema.ForOptions{
		TypeSchemas: FalseTypeSchemaTypes,
	})
	if err != nil {
		panic(err)
	}
	return s
})
