package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

func (a *Application) HandleSignup(w http.ResponseWriter, r *http.Request) {
	// Read Input.
	var userRequest storage.UserRequest

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
	userJWT, err := a.DB.User.CreateUser(&userRequest)

	if err != nil {
		switch err {
		case storage.ErrNoRecord:
			a.BadRequest(w)
		case storage.ErrDuplicateEntry:
			a.ClientError(w, http.StatusConflict)
		default:
			a.ServerError(w, err)
		}
		return
	}

	// produce token
	token, err := IssueUserJWT(*userJWT, a.Config.AccessToken)
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
	var userRequest storage.UserRequest

	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// validate Request
	if err := a.Validate.Struct(userRequest); err != nil {
		a.BadRequest(w)
		return
	}

	userJWT, err := a.DB.User.CheckUserExists(&userRequest)
	if err != nil {
		if err == storage.ErrInvalidCredentials {
			a.ClientError(w, http.StatusUnauthorized)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	// produce token
	token, err := IssueUserJWT(*userJWT, a.Config.AccessToken)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	// send token

	Response := &struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Message: fmt.Sprint("Hello ", userRequest.UserName),
		Token:   token,
	}

	WriteJSON(w, http.StatusAccepted, Response)
}

func (a *Application) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	user, err := a.DB.User.GetUser(0)
	if err != nil {
		if err == storage.ErrNoRecord {
			a.NotFound(w)
		} else {
			a.ServerError(w, err)
		}
		return
	}

	WriteJSON(w, http.StatusOK, user)
}
