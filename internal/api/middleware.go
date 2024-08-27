package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

// Logger used to log requests. For each session it prints
// http Protocol, http Method, url path
func (a *Application) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.InfoLog.Printf("%v: %v %v\n", r.Proto, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func (a *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Create a deferred function (which will always be run in the event
		// of a panic as Go unwinds the stack).
		defer func() {
			// Use the builtin recover function to check if there has been a
			// panic or not. If there has...
			if err := recover(); err != nil {
				// Set a "Connection: close" header on the response.
				w.Header().Set("Connection", "close")
				// Call the app.serverError helper method to return a 500
				// Internal Server response.
				a.ServerError(w, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", " default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")
		next.ServeHTTP(w, r)
	})
}

// Authenticate Verifies the JWT Token Signature.
func (a *Application) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the JWT token from the request header 'Authorization'
		authorizationHeader := r.Header.Get("Authorization")
		tokenString := strings.TrimPrefix(authorizationHeader, "Bearer ")

		// Check if the token has the prefix Bearer or not
		if authorizationHeader == tokenString {
			a.ClientError(w, http.StatusUnauthorized)
			return
		}

		// Verifies the claims
		claims, err := VerifyUserJWT(tokenString, a.Config.AccessToken)
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

// AuthorizeUserInfo Check whether the User is authorized to access
// the resources
func (a *Application) AuthorizeUserInfo(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userTarget := r.PathValue("user")

		userJWT := r.Context().Value("userJWT").(*storage.UserJWT)

		if userJWT.UserName != userTarget && userJWT.Role != "admin" {
			fmt.Printf("userJWT.UserName: %v\n", userJWT.UserName)
			fmt.Printf("userTarget: %v\n", userTarget)
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
func (a *Application) AuthorizeAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userJWT := r.Context().Value("userJWT").(*storage.UserJWT)

		if userJWT.Role != "admin" {
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

func (a *Application) RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := strings.Split(r.RemoteAddr, "]")[0] + "]"

		ipCacheKey := fmt.Sprintf("rateLimit : %v", ip)

		attempts, err := a.Cache.Incr(ipCacheKey)
		if err != nil {
			a.ServerError(w, err)
			return
		}

		if attempts == 1 {
			err := a.Cache.Expire(ipCacheKey, a.Config.RateDuration)
			if err != nil {
				a.ServerError(w, err)
				return
			}
		} else if attempts >= a.Config.RateLimit {
			a.ClientError(w, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
