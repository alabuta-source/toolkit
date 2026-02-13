package apiError

import (
	"fmt"
)

type RequestError interface {
	Code() string
	ErrorCause() string
	Error() string
	Status() int
	Message() string
}

type requestError struct {
	ErrorMessage string `json:"message"`
	ErrorCode    string `json:"error,omitempty"`
	Cause        string `json:"cause,omitempty"`
	Causes       []any  `json:"causes,omitempty"`
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

func (e requestError) ErrorCause() string {
	return e.Cause
}
