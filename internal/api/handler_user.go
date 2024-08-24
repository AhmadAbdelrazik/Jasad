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

// HandleSignup Handles the signup operation
// It checks if the username is available, the password is strong enough.
// Creates new user and send Token for user.
func (a *Application) HandleSignup(w http.ResponseWriter, r *http.Request) {
	// Initialize UserRequest object to parse the request body
	var userRequest storage.UserRequest

	// Parse the request body using json, returns error at failure
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// Validate the response
	if err := a.Validate.Struct(userRequest); err != nil {
		a.BadRequest(w)
		return
	}

	// Database Call
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

	// Issue User JWT Token
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

// HandleSignIn Handles the signin operations
// It checks if the username and password are valid.
// Count the login attempts on valid usernames.
func (a *Application) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	// Initialize UserRequest object to parse the request body
	var userRequest storage.UserRequest

	// Parse the request body using json, returns error at failure
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		a.BadRequest(w)
		return
	}

	defer r.Body.Close()

	// Validate the response
	if err := a.Validate.Struct(userRequest); err != nil {
		a.BadRequest(w)
		return
	}

	// Check login attempts in cache
	var loginAttempts int

	// The Cache hit returns a string, we need to convert it
	// to int to track login attempts
	attemptsRaw, err := a.Cache.Get(fmt.Sprintf("username: %s", userRequest.UserName))
	if err != nil {
		if err == cache.ErrNotExist {
			attemptsRaw = "0"
		} else {
			a.ServerError(w, err)
			return
		}
	}

	// Convert the Login attempts to int
	loginAttempts, err = strconv.Atoi(attemptsRaw)
	if err != nil {
		a.ServerError(w, err)
		return
	}

	// Check login attempts
	if loginAttempts >= a.Config.LoginAttemptsLimit {
		a.ClientError(w, http.StatusTooManyRequests)
		return
	}

	// Call for the Database
	userJWT, err := a.DB.User.CheckUserExists(&userRequest)
	if err != nil {
		switch err {
		case storage.ErrInvalidUsername:
			a.ClientError(w, http.StatusUnauthorized)

		// Register login attempt on valid usernames only.
		// This Help keep the Cache clean from invalid usernames, and
		// record login attempts on existing usernames only.
		case storage.ErrInvalidPassword:
			err := a.Cache.Set(
				fmt.Sprintf("username: %s", userRequest.UserName),
				fmt.Sprintf("%d", loginAttempts+1),
				5*time.Minute)

			if err != nil {
				a.ServerError(w, err)
				return
			}
			a.ClientError(w, http.StatusUnauthorized)
		default:
			a.ServerError(w, err)
		}
		return
	}

	// Issue User JWT Token
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

// HandleUserInfo Returns The userinfo containing username, role, and userID.
// WARNING: This Handlers should be accessed only by the ADMIN users, as userIDs
// are SENSETIVE INFORMATION.
func (a *Application) HandleUserInfo(w http.ResponseWriter, r *http.Request) {
	// Get the UserID (Parsed Earlier in the auth middleware)
	userID := r.Context().Value("userID").(int)

	// Database Call
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
