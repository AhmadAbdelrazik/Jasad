package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
)

func (s *APIServer) HandleSignup(w http.ResponseWriter, r *http.Request) error {
	// Read Input.
	userRequest := storage.UserCreateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return s.BadRequest()
	}

	// Validate Struct
	if err := s.Validate.Struct(userRequest); err != nil {
		return s.BadRequest()
	}

	// Add user to database
	userID, err := s.DB.CreateUser(&userRequest)

	if err != nil {
		if err == storage.ErrNoRecord {
			return s.BadRequest()
		} else if strings.Contains(err.Error(), "Duplicate entry") {
			return s.ClientError(http.StatusConflict)
		} else {
			return s.ServerError(err)
		}
	}

	// produce token
	token, err := IssueUserJWT(userID.String(), "user")
	if err != nil {
		return s.ServerError(err)
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

	return nil
}

func (s *APIServer) HandleSignIn(w http.ResponseWriter, r *http.Request) error {
	// Read Input
	var userSigninRequest storage.UserSigninRequest

	if err := json.NewDecoder(r.Body).Decode(&userSigninRequest); err != nil {
		return s.BadRequest()
	}

	// validate Request
	if err := s.Validate.Struct(userSigninRequest); err != nil {
		return s.BadRequest()
	}

	userID, err := s.DB.CheckUserExists(&userSigninRequest)
	if err != nil {
		if err == storage.ErrNoRecord {
			return s.NotFound()
		} else if err == storage.ErrInvalidCredentials {
			return s.ClientError(http.StatusUnauthorized)
		} else {
			return s.ServerError(err)
		}
	}

	// produce token
	token, err := IssueUserJWT(userID.String(), "user")
	if err != nil {
		return s.ServerError(err)
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

	return nil
}
