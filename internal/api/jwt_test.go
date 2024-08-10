package api

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestJWT(t *testing.T) {
	testSecret := []byte(`my-test-jwt-secret`)
	t.Run("Sign a token", func(t *testing.T) {
		claims := UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "test subject",
				ID:      "001",
			},
		}

		got, err := SignJWT(claims, testSecret)
		want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0IHN1YmplY3QiLCJqdGkiOiIwMDEifQ.K97_5j2YcWOzT_7MSSYisJ15IFNbKbax5TsEvwvKMJI"
		assertEquality(t, got, want)
		assertNoError(t, err)
	})
}

func assertEquality(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot %v\nwant %v", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got %v want nil", err)
	}
}
