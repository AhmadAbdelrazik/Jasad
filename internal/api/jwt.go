package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SignJWT(claims jwt.MapClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// For Security Reasons, you should NEVER hard code
	// sensitive information like JWT key and instead store it
	// in environment variables and call it using os.Getenv("key-name")
	key := "69zDfhhZUxnNl63VqmV3EQWja9++RsqORbltMyeTMVHm"

	return t.SignedString([]byte(key))
}

func IssueUserJWT(userID, role string) (string, error) {
	claims := jwt.MapClaims{
		"exp":    time.Now().Add(15 * time.Minute),
		"userID": userID,
		"role":   role,
	}

	return SignJWT(claims)
}
