package main

import "net/http"

type ErrBadRequest struct {
	Status  int
	Message string
}

func NewErrBadRequest() *ErrBadRequest {
	return &ErrBadRequest{
		Message: "Bad Request",
		Status: http.StatusBadRequest,
	}
}

func (r ErrBadRequest) Error() string {
	return r.Message
}

type ErrNotFound struct {
	Status  int
	Message string
}

func NewErrNotFound(message string) *ErrNotFound {
	return &ErrNotFound{
		Message: message,
		Status: http.StatusNotFound,
	}
}

func (r ErrNotFound) Error() string {
	return r.Message
}