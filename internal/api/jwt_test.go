package api

import (
	"errors"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/go-cmp/cmp"
)

func TestJWT(t *testing.T) {
	t.Run("Sign a token", func(t *testing.T) {
		testSecret := []byte(`my-test-jwt-secret`)
		claims := &UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "test subject",
				ID:      "001",
			},
		}

		got, err := SignJWT(claims, testSecret)
		want := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0IHN1YmplY3QiLCJqdGkiOiIwMDEifQ.K97_5j2YcWOzT_7MSSYisJ15IFNbKbax5TsEvwvKMJI"
		assertStringEquality(t, got, want)
		assertNoError(t, err)
	})

	t.Run("Verify a token", func(t *testing.T) {
		claims := &UserClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: "test subject",
				ID:      "002",
			},
		}
		mySecret := `correct-secret`
		notMySecret := `incorrect-secret`
		tokenString, _ := SignJWT(claims, []byte(mySecret))
		t.Run("Verifying correct token", func(t *testing.T) {
			got, err := VerifyJWT(tokenString, mySecret)
			assertClaimEquality(t, got, claims)
			assertNoError(t, err)
		})

		t.Run("Verifying incorrect token", func(t *testing.T) {
			got, err := VerifyJWT(tokenString, notMySecret)
			assertNoClaim(t, got)
			assertError(t, err, jwt.ErrInvalidKey)
		})
	})
}

func assertStringEquality(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("\ngot %v\nwant %v", got, want)
	}
}

func assertError(t testing.TB, got, want error) {
	t.Helper()
	if errors.Is(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertNoError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("got %v want nil", err)
	}
}

func assertNoClaim(t testing.TB, got *UserClaims) {
	t.Helper()
	if got != nil {
		t.Errorf("got %v want nil", got)
	}
}

func assertClaimEquality(t testing.TB, got, want *UserClaims) {
	t.Helper()
	if !cmp.Equal(got, want) {
		t.Errorf("\ngot %v\nwant %v", got, want)
	}
}
