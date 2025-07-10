package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

// Reset the user's password. Authentication is done with the JWT.
func userResetPassword(w http.ResponseWriter, r *http.Request) {
	_, userID, password, err := checkUserRequest[string](r)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	newPasswordHash, err := createPasswordHash(*password)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	database.SetUserPassword(userID, newPasswordHash)
	Coms.ExternalPostRespondCode(http.StatusOK, w)
}
