package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

func userRequestData(w http.ResponseWriter, r *http.Request) {
	_, userID, _, err := checkUserRequest[any](r)
	if err != nil {
		Coms.PrintErrStr("Could not get user data: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	user, err := database.GetUserData(userID)
	if err != nil {
		Coms.PrintErrStr("Could not get user data: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	Coms.ExternalPostRespond(user, w)
}
