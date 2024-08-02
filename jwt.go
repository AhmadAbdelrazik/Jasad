package main

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignJWT(claims jwt.MapClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// For Security Reasons, you should NEVER hard code
	// sensitive information like JWT key and instead store it
	// in environment variables and call it using os.Getenv("key-name")
	key := "69zDfhhZUxnNl63VqmV3EQWja9++RsqORbltMyeTMVHm"

	return t.SignedString(key)
}

func IssueUserJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"iat":    time.Now().UTC(),
		"userID": userID,
		"role":   role,
	}

	return SignJWT(claims)
}
