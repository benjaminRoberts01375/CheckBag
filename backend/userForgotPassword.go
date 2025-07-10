package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/email"
)

func userForgotPasswordRequest(w http.ResponseWriter, r *http.Request) {
	userEmail, err := Coms.ExternalPostReceived[string](r)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	transactionID, err := cache.setForgotPassword(*userEmail)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}

	go email.SendForgotPassword(*userEmail, transactionID)

	Coms.ExternalPostRespondCode(http.StatusOK, w)
}

func userForgotPasswordCheckValid(w http.ResponseWriter, r *http.Request) {
	replaceToken := r.PathValue("token")
	email, err := cache.getForgotPassword(replaceToken)

	if err != nil || email == "" {
		Coms.ExternalPostRespondCode(http.StatusNotFound, w)
		return
	}
	Coms.ExternalPostRespondCode(http.StatusOK, w)
}

func userForgotPasswordConfirmation(w http.ResponseWriter, r *http.Request) {
	replaceToken := r.PathValue("token")
	newPassword, err := Coms.ExternalPostReceived[string](r)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	email, err := cache.getAndDeleteResetPassword(replaceToken)
	if err != nil || email == "" {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	newPasswordHash, err := createPasswordHash(*newPassword)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	database.SetUserForgotPassword(email, newPasswordHash)
	Coms.ExternalPostRespondCode(http.StatusOK, w)
}
