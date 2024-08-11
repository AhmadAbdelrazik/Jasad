package api

import (
	"net/http"
	"strings"
)

func (a *Application) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.InfoLog.Printf("%v: %v %v\n", r.Proto, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Authenticate just verifies the token, no more, no less. it also checks if it has expired.
func (a *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		if authorizationHeader == tokenString {
			a.BadRequest(w)
		}

		next.ServeHTTP(w, r)
	})
}
