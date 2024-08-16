package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

func (a *Application) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.InfoLog.Printf("%v: %v %v\n", r.Proto, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// Authenticate Verifies the JWT Token Signature.
func (a *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		if authorizationHeader == tokenString {
			a.ClientError(w, http.StatusUnauthorized)
			return
		}

		claims, err := VerifyJWT(tokenString, a.Config.AccessToken)
		if err != nil {
			a.ClientError(w, http.StatusUnauthorized)
			return
		}

		userJWT := &storage.UserJWT{
			UserName: claims.Subject,
			Role:     claims.Role,
		}

		ctx := context.WithValue(r.Context(), "userJWT", userJWT)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *Application) AuthorizeUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userTarget := r.PathValue("user")

		userJWT := r.Context().Value("userJWT").(storage.UserJWT)

		if userJWT.UserName != userTarget && userJWT.Role != "admin" {
			a.ClientError(w, http.StatusUnauthorized)
			return
		}

		userID, err := a.DB.User.GetUserID(userJWT.UserName)
		if err != nil {
			if err == storage.ErrNoRecord {
				a.NotFound(w)
			} else {
				a.ServerError(w, err)
			}

			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
