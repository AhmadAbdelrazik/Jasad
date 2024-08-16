package api

import (
	"fmt"
	"time"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func SignJWT(claims *UserClaims, secret []byte) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString(secret)
}

func IssueUserJWT(user storage.UserJWT, secret string) (string, error) {
	claims := &UserClaims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserName,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute).UTC()),
		},
	}

	fmt.Printf("claims: %v\n", claims)
	return SignJWT(claims, []byte(secret))
}

func VerifyJWT(tokenString string, secret string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*UserClaims), nil
}
