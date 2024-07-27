package main

import "net/http"

func (s *APIServer)logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.InfoLog.Printf("%v: %v %v\n", r.Proto, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}