package application

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/ahmadabdelrazik/jasad/internal/model"
	"golang.org/x/time/rate"
)

func (app *Application) IsAuthorized(next http.HandlerFunc, accpetedRoles ...model.Role) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token
		cookie, err := r.Cookie("id")
		if err != nil { // the only possible error is http.ErrNoCookie
			AuthenticationErrorResponse(w, r)
			return
		}

		token := cookie.Value
		// get the corresponding session
		session, err := app.models.Tokens.GetSessionFromToken(token)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrNotFound):
				AuthenticationErrorResponse(w, r)
			default:
				ServerErrorResponse(w, r, err)
			}
			return
		}

		accpetedRoles = append(accpetedRoles, model.RoleAdmin)

		// authorize
		if !slices.Contains(accpetedRoles, session.Role) {
			UnauthorizedResponse(w, r)
			return
		}

		// load the full user from database
		user, err := app.models.Users.GetByID(session.UserID)
		if err != nil {
			ServerErrorResponse(w, r, err)
		}

		// place it in the request to be fetched by the handlers
		r = withUser(r, user)

		next.ServeHTTP(w, r)
	}
}

func (app *Application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "Close")

				ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *Application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.cfg.LimiterEnable {
			next.ServeHTTP(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ServerErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			clients[ip] = &client{limiter: rate.NewLimiter(
				rate.Limit(app.cfg.LimiterRPS),
				app.cfg.LimiterBurst,
			)}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			RateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}

type contextkey string

const userContext contextkey = "user"

func withUser(r *http.Request, user *model.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContext, user)
	return r.WithContext(ctx)
}

func getUser(r *http.Request) (*model.User, bool) {
	user, ok := r.Context().Value(userContext).(*model.User)
	return user, ok
}
