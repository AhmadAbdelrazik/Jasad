package api

import (
	"context"
	"net/http"
	"strings"
)

func (a *Application) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.InfoLog.Printf("%v: %v %v\n", r.Proto, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (a *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		if authorizationHeader == tokenString {
			a.BadRequest(w)
			return
		}

		claims, err := VerifyJWT(tokenString, []byte(a.Config.AccessToken))
		if err != nil {
			a.ClientError(w, http.StatusUnauthorized)
			return
		}

		targetUser := r.PathValue("user")

		if claims.Username != targetUser {
			a.ClientError(w, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "jwt", claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
