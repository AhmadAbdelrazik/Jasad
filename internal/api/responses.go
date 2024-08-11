package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

func (a *Application) ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	a.ErrorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *Application) ClientError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}

func (a *Application) BadRequest(w http.ResponseWriter) {
	a.ClientError(w, http.StatusBadRequest)
}
func (a *Application) NotFound(w http.ResponseWriter) {
	a.ClientError(w, http.StatusNotFound)
}
