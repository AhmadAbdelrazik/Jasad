package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (app *APIServer) ServerError(err error) *apiErr {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	app.ErrorLog.Output(2, trace)
	return &apiErr{Message: http.StatusText(http.StatusInternalServerError), Status: http.StatusInternalServerError}
}

func (app *APIServer) ClientError(code int) *apiErr {
	return &apiErr{Message: http.StatusText(code), Status: code}
}

func (app *APIServer) BadRequest() *apiErr {
	return app.ClientError(http.StatusBadRequest)
}
func (app *APIServer) NotFound() *apiErr {
	return app.ClientError(http.StatusNotFound)
}
