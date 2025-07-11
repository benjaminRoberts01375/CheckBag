package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

func userExists(w http.ResponseWriter, r *http.Request) {
	data, err := fileSystem.GetUserData()
	if err != nil || data == "" {
		Coms.Println("User does not exist")
		Coms.ExternalPostRespondCode(http.StatusGone, w)
		return
	}
	Coms.Println("User exists")
	Coms.ExternalPostRespondCode(http.StatusOK, w)
}
