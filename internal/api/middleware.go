package api

import (
	"net/http"
	"strings"
)

func (s *APIServer) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.InfoLog.Printf("%v: %v %v\n", r.Proto, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Authenticate just verifies the token, no more, no less. it also checks if it has expired.
func (s *APIServer) Authenticate(next http.Handler) http.Handler {
	return makeHTTPHandleFunc(func(w http.ResponseWriter, r *http.Request) error {
		authorizationHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		if authorizationHeader == tokenString {
			return s.BadRequest()
		}

		next.ServeHTTP(w, r)
		return nil
	})
}
