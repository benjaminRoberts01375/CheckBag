package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/email"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/models"
	"golang.org/x/crypto/bcrypt"
)

func newUser(w http.ResponseWriter, r *http.Request) {
	userData, err := Coms.ExternalPostReceived[models.UserCreate](r)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	userPassword, err := createPasswordHash(userData.Password)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	userData.Password = string(userPassword)
	database.SetNewUser(*userData)

	Coms.ExternalPostRespondCode(http.StatusOK, w)
	activationToken, err := cache.setNewUser(userData.Email)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	go email.SendNewUser(userData.Email, userData.FirstName, activationToken)
}

func createPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10) // 72 bytes max
}

func newUserConfirmation(w http.ResponseWriter, r *http.Request) {
	activationToken := r.PathValue("token")
	username, err := cache.getAndDeleteNewUser(activationToken)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	http.Redirect(w, r, "https://"+email.Config.BaseURL+"/login", http.StatusSeeOther)
	database.SetNewUserConfirmed(username)
}
