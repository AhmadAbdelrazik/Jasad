package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

func (a *Application) HandleSignup(w http.ResponseWriter, r *http.Request) {
	// Read Input.
	userRequest := storage.UserCreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// Validate Struct
	if err := a.Validate.Struct(userRequest); err != nil {
		a.BadRequest(w)
		return
	}

	// Add user to database
	userID, err := a.DB.CreateUser(&userRequest)

	if err != nil {
		if err == storage.ErrNoRecord {
			a.BadRequest(w)
		} else if strings.Contains(err.Error(), "Duplicate entry") {
			a.ClientError(w, http.StatusConflict)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	// produce token
	token, err := IssueUserJWT(userID.String(), "user", []byte(a.Config.AccessToken))
	if err != nil {
		a.ServerError(w, err)
		return
	}

	// send token

	Response := &struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: "user has been created",
		Token:   token,
	}

	WriteJSON(w, http.StatusCreated, Response)

}

func (a *Application) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	// Read Input
	var userSigninRequest storage.UserSigninRequest

	if err := json.NewDecoder(r.Body).Decode(&userSigninRequest); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// validate Request
	if err := a.Validate.Struct(userSigninRequest); err != nil {
		a.BadRequest(w)
		return
	}

	userID, err := a.DB.CheckUserExists(&userSigninRequest)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else if err == storage.ErrInvalidCredentials {
			a.ClientError(w, http.StatusUnauthorized)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	// produce token
	token, err := IssueUserJWT(userID.String(), "user", []byte(a.Config.AccessToken))
	if err != nil {
		a.ServerError(w, err)
		return
	}

	// send token

	Response := &struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: fmt.Sprint("Hello ", userSigninRequest.UserName),
		Token:   token,
	}

	WriteJSON(w, http.StatusCreated, Response)
}
