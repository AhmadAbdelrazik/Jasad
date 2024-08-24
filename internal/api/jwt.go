package api

import (
	"time"

	"github.com/AhmadAbdelrazik/jasad/internal/storage"
	"github.com/golang-jwt/jwt/v5"
)

// UserClaims These are the claims in the JWT.
// We use the following fields:
//
// Role: select the user role 'admin', 'user'
// Subject: Contains the username.
// IssuedAt: Contains the Issuing date and time.
// ExpiresAt: Contains the expiration date and time.
type UserClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// SignJWT Takes claims and secret and return the token.
// It uses HS256 as the main encryption method.
func SignJWT(claims *UserClaims, secret []byte) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return t.SignedString(secret)
}

// IssueUserJWT Issues a new token for a user.
func IssueUserJWT(user storage.UserJWT, secret string) (string, error) {
	claims := &UserClaims{
		Role: user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.UserName,
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute).UTC()),
		},
	}

	return SignJWT(claims, []byte(secret))
}

// VerifyUserJWT Verifies the authenticity of the claim, returns UserClaims at success
func VerifyUserJWT(tokenString string, secret string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	return token.Claims.(*UserClaims), nil
}
