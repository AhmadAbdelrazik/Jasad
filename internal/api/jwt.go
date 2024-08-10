package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	jwt.RegisteredClaims
}

// For Security Reasons, you should NEVER hard code
// sensitive information like JWT key and instead store it
// in environment variables and call it using os.Getenv("key-name")
var key = []byte("69zDfhhZUxnNl63VqmV3EQWja9++RsqORbltMyeTMVHm")

func SignJWT(claims UserClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString(key)
}

func IssueUserJWT(userID, role string) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        userID,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute).UTC()),
		},
	}

	return SignJWT(claims)
}

func VerifyJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*UserClaims), nil
}
