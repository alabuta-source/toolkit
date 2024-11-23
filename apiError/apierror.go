package apiError

import (
	"fmt"
)

type RequestError interface {
	Code() string
	Error() string
	Status() int
	Message() string
}

type requestError struct {
	ErrorMessage string         `json:"message"`
	ErrorCode    string         `json:"error"`
	Cause        string         `json:"cause"`
	Causes       map[string]any `json:"causes"`
	ErrorStatus  int            `json:"status"`
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
