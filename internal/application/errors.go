package application

import (
	"encoding/json"
	"maps"
	"net/http"

	"github.com/rs/zerolog/log"
)

func logError(r *http.Request, err error) {
	log.Error().
		Stack().Err(err).
		Str("request_method", r.Method).
		Str("request_url", r.URL.String()).
		Msg("")
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"

	ErrorResponse(w, r, http.StatusNotFound, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func ConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "resource already exists"
	ErrorResponse(w, r, http.StatusConflict, message)
}

func EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to edit conflict, please try again"
	ErrorResponse(w, r, http.StatusConflict, message)
}

func AuthenticationErrorResponse(w http.ResponseWriter, r *http.Request) {
	message := "Invalid Authentication Credentials"
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func UnauthorizedResponse(w http.ResponseWriter, r *http.Request) {
	message := "Insufficient permission to access the resource"
	ErrorResponse(w, r, http.StatusForbidden, message)
}

func RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"

	ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	maps.Copy(w.Header(), headers)

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)
	w.Write(js)

	return nil
}
