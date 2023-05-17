package apierr

import (
	"fmt"
	"net/http"
)

var (
	ErrInvalidRequest      = New(http.StatusBadRequest, "InvalidRequest", "Request form is not valid.")
	ErrResourceNotFound    = New(http.StatusNotFound, "ResourceNotFound", "Resource not found.")
	ErrAuthenticationFail  = New(http.StatusUnauthorized, "FailedAuthentication", "Auth failed")
	ErrResourceConflict    = New(http.StatusConflict, "ResourceAlreadyExist", "Resource already exists")
	ErrInternalServerError = New(http.StatusInternalServerError, "InternalServerError",
		"There was an error. Please try again later.")
)

type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
	RequestID  string `json:"requestId,omitempty"`
}

func New(statusCode int, code string, msg string) *Error {
	return &Error{StatusCode: statusCode, Code: code, Message: msg}
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) WithStatusCode(statusCode int) *Error {
	return &Error{StatusCode: statusCode, Code: e.Code, Message: e.Message}
}

func (e *Error) WithCode(code string) *Error {
	return &Error{StatusCode: e.StatusCode, Code: code, Message: e.Message}
}

func (e *Error) WithMessage(msg string) *Error {
	return &Error{StatusCode: e.StatusCode, Code: e.Code, Message: msg}
}

func (e *Error) WithMessagef(format string, args ...any) *Error {
	return &Error{StatusCode: e.StatusCode, Code: e.Code, Message: fmt.Sprintf(format, args...)}
}
