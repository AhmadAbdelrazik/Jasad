package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// ServerError Log the stack trace, and send http error with code
// InternalServerError
func (a *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	a.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// ClientError writes client error to the response, with the default
// status text
func (a *Application) ClientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

// BadRequest writes bad request error to the response
func (a *Application) BadRequest(w http.ResponseWriter) {
	a.ClientError(w, http.StatusBadRequest)
}

// NotFound writes not found error to the response
func (a *Application) NotFound(w http.ResponseWriter) {
	a.ClientError(w, http.StatusNotFound)
}
