package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/email"
)

func userResetEmailRequest(w http.ResponseWriter, r *http.Request) {
	claims, _, updatedEmail, err := checkUserRequest[string](r)
	if err != nil {
		Coms.PrintErrStr("Checking user request in reset email request: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	activationToken, err := cache.setChangeEmail(claims.Username, *updatedEmail)
	if err != nil {
		Coms.PrintErrStr("Cache error during reset email request: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	go email.SendResetEmail(*updatedEmail, activationToken)
	Coms.ExternalPostRespondCode(http.StatusOK, w)
}

func userResetEmailConfirmation(w http.ResponseWriter, r *http.Request) {
	activationToken := r.PathValue("token")
	oldEmail, newEmail, err := cache.getAndDeleteChangeEmail(activationToken)
	if err != nil {
		Coms.PrintErrStr("Cache error during reset email confirmation: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}

	err = database.SetUserEmail(oldEmail, newEmail)
	if err != nil {
		Coms.PrintErrStr("Error during reset email confirmation: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	http.Redirect(w, r, "https://"+email.Config.BaseURL+"/login", http.StatusSeeOther)
}
