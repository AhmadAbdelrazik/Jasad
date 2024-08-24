package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/AhmadAbdelrazik/jasad/internal/cache"
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

	// Check login attempts in cache
	var loginAttempts int
	attemptsRaw, err := a.Cache.Get(fmt.Sprintf("username: %s", userRequest.UserName))
	if err != nil {
		if err == cache.ErrNotExist {
			attemptsRaw = "0"
		} else {
			a.ServerError(w, err)
			return
		}
	}

	loginAttempts, err = strconv.Atoi(attemptsRaw)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	fmt.Printf("loginAttempts: %v\n", loginAttempts)

	if loginAttempts >= 5 {
		a.ClientError(w, http.StatusTooManyRequests)
		return
	}

	// Call for the Database
	userJWT, err := a.DB.User.CheckUserExists(&userRequest)
	if err != nil {
		if err == storage.ErrInvalidUsername {
			a.ClientError(w, http.StatusUnauthorized)
		} else if err == storage.ErrInvalidPassword {
			// Register login attempt on the username.
			// We don't register globally to protect the cache from filling the heap
			// with invalid usernames, and only set our cache on actual usernames
			if err := a.Cache.Set(
				fmt.Sprintf("username: %s", userRequest.UserName),
				fmt.Sprintf("%d", loginAttempts+1),
				5*time.Minute); err != nil {
				a.ServerError(w, err)
				return
			}

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
	userID := r.Context().Value("userID").(int)
	user, err := a.DB.User.GetUser(userID)
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
