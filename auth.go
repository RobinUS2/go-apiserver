package apiserver

import (
	"github.com/dgrijalva/jwt-go"
)

type ApiAuth struct {
	hmacSecret []byte
}

func (auth *ApiAuth) NewTestToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"testing": true,
		// "foo": "bar",
		// "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(auth.hmacSecret)
	return tokenString, err
}

func NewAuth() *ApiAuth {
	return &ApiAuth{
		hmacSecret: []byte("PEk7Dg90VJqEWGnbnMDazHZLY9r08UaOH7AI4O9UfJbF2n9Uq0y25k81xpUFmhoT"),
	}
}
