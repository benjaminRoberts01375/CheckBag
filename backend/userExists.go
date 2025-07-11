package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

func userExists(w http.ResponseWriter, r *http.Request) {
	data, err := fileSystem.GetUserData()
	if err != nil || data == "" {
		Coms.ExternalPostRespondCode(http.StatusNotFound, w)
		return
	}
	Coms.ExternalPostRespondCode(http.StatusOK, w)
}
