package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ahmadabdelrazik/jasad/internal/model"
	"github.com/ahmadabdelrazik/jasad/pkg/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type OAuthConfig struct {
	Google oauth2.Config
}

func newOAuthConfig(cfg config.Config) OAuthConfig {
	google := oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  fmt.Sprintf("%s/google_callback", cfg.Origin),
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	return OAuthConfig{
		Google: google,
	}
}

// InfoToken contains the information received from google.
type InfoToken struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// googleLoginHandler redirect user to google's login page to consent using the
// application
func (app *Application) googleLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := app.oauth.Google.AuthCodeURL("random-state-to-protect-from-csrf")

	http.Redirect(w, r, url, http.StatusSeeOther)
}

// googleCallbackHandler receive auth code and exchange it with user info from
// google auth server.
func (app *Application) googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query()

	state := app.readString(qs, "state", "")
	if state != "random-state-to-protect-from-csrf" {
		AuthenticationErrorResponse(w, r)
		return
	}

	code := app.readString(qs, "code", "")

	// exchange user code with authorized token from google.
	token, err := app.oauth.Google.Exchange(context.Background(), code)
	if err != nil {
		AuthenticationErrorResponse(w, r)
		return
	}

	// exchange the auth token from google with the information needed.
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		AuthenticationErrorResponse(w, r)
		return
	}

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		ServerErrorResponse(w, r, fmt.Errorf("failed to read response body: %w", err))
		return
	}
	defer resp.Body.Close()

	var info InfoToken

	if err := json.Unmarshal(userData, &info); err != nil {
		ServerErrorResponse(w, r, fmt.Errorf("failed to unmarshal userData: %w", err))
		return
	}

	user, err := fetchUser(app.models, info)
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}

	// Produce Session Token to give to the user.
	sessionToken, err := app.models.Tokens.GenerateToken(user)
	if err != nil {
		ServerErrorResponse(w, r, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:       "id",
		Value:      sessionToken,
		Expires:    time.Now().Add(3 * 24 * time.Hour),
		RawExpires: "",
		MaxAge:     0,
		Secure:     true,
		HttpOnly:   true,
		SameSite:   http.SameSiteLaxMode,
	})

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "logged in successfully"}, nil)
	if err != nil {
		ServerErrorResponse(w, r, err)
	}
}

// fetchUser returns the user by email, or create a new one if not
// registered
func fetchUser(models *model.Model, info InfoToken) (*model.User, error) {
	user, err := models.Users.GetByEmail(info.Email)
	if user != nil { // user is registered on the system
		return user, nil
	} else if errors.Is(err, model.ErrNotFound) { // user is not registered

		user := &model.User{
			Name:  info.Name,
			Email: info.Email,
			Role:  model.RoleUser,
		}

		if err := models.Users.Create(user); err != nil {
			return nil, err
		}

		return user, nil
	} else { // unexpected error
		return nil, err
	}
}
