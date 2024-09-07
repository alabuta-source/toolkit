package apiError

import (
	"fmt"
)

type ErrorCauses []any

type RequestError interface {
	ErrorCode() string
	Error() string
	ErrorStatus() int
	ErrorMessage() string
}

type requestError struct {
	Message string      `json:"message"`
	Code    string      `json:"error"`
	Status  int         `json:"status"`
	Cause   ErrorCauses `json:"cause"`
}

func (e requestError) ErrorCode() string {
	return e.Code
}

func (e requestError) Error() string {
	return fmt.Sprintf("Message: %s;Error Code: %s;Status: %d; Cause: %v", e.Message, e.Code, e.Status, e.Cause)
}

func (e requestError) ErrorStatus() int {
	return e.Status
}

func (e requestError) ErrorMessage() string {
	return e.Message
}
