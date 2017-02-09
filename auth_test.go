package apiserver_test

import (
	"log"
	"testing"

	"../apiserver"
)

func TestAuth(t *testing.T) {
	auth := apiserver.NewAuth()
	res, err := auth.NewTestToken()
	log.Printf("%v %v", res, err)
}
