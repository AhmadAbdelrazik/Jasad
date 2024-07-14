package main

import "net/http"

func (app *Application) logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}