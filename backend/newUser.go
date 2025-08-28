package main

import (
	"net/http"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
	"golang.org/x/crypto/bcrypt"
)

func newUser(w http.ResponseWriter, r *http.Request) {
	userData, _ := fileSystem.GetUserData()
	if userData != "" {
		requestRespondCode(w, http.StatusForbidden)
		return
	}

	newUserPassword, err := requestReceived[string](r)
	if err != nil {
		requestRespondCode(w, http.StatusBadRequest)
		return
	}
	userPasswordHash, err := createPasswordHash(*newUserPassword)
	if err != nil {
		requestRespondCode(w, http.StatusBadRequest)
		return
	}
	err = fileSystem.SetUserData(string(userPasswordHash))
	if err != nil {
		Printing.PrintErrStr("Could not set user data in file system: ", err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}
	setJWTCookie(w)
	requestRespondCode(w, http.StatusOK)
}

func createPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10) // 72 bytes max
}
