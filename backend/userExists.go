package main

import (
	"net/http"

	Printing "github.com/benjaminRoberts01375/Web-Tech-Stack/logging"
)

func userExists(w http.ResponseWriter, r *http.Request) {
	data, err := fileSystem.GetUserData()
	if err != nil || data == "" {
		Printing.Println("User does not exist")
		requestRespondCode(w, http.StatusGone)
		return
	}
	Printing.Println("User exists")
	requestRespondCode(w, http.StatusOK)
}
