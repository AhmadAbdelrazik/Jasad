package application

import (
	"context"
	"errors"
	"net/http"
	"slices"

	"github.com/ahmadabdelrazik/jasad/internal/model"
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
