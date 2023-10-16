package apiError

import (
	"errors"
	"fmt"
	"net/http"
)

type RequestError interface {
	Message() string
	Status() int
	Error() string
}

type RequestEr struct {
	ErrorMessage string `json:"message"`
	ErrorCode    string `json:"error"`
	ErrorStatus  int    `json:"status"`
}

type requestError struct {
	ErrorMessage string `json:"message"`
	ErrorCode    string `json:"error"`
	ErrorStatus  int    `json:"status"`
}

func (e requestError) Code() string {
	return e.ErrorCode
}

func (e requestError) Error() string {
	return fmt.Sprintf("Message: %s;Error Code: %s;Status: %d", e.ErrorMessage, e.ErrorCode, e.ErrorStatus)
}

func (e requestError) Status() int {
	return e.ErrorStatus
}

func (e requestError) Message() string {
	return e.ErrorMessage
}

func NewRequestError(message string, error string, status int) RequestError {
	return requestError{message, error, status}
}

func NewForbiddenError(message string) RequestError {
	return requestError{ErrorMessage: message, ErrorCode: "forbidden", ErrorStatus: http.StatusForbidden}
}

func NewBadRequestError(message string) RequestError {
	return requestError{ErrorMessage: message, ErrorCode: "bad_request", ErrorStatus: http.StatusBadRequest}
}

func NewUnauthorizedError(message string) RequestError {
	return requestError{message, "unauthorized", http.StatusUnauthorized}
}

func NewInternalServerApiError(message string, err error) RequestError {
	return requestError{message, err.Error(), http.StatusInternalServerError}
}

func ParseError(err error) RequestError {
	var apiErr RequestError
	ok := errors.As(err, &apiErr)
	if !ok {
		return NewInternalServerApiError(err.Error(), err)
	}
	return apiErr
}
