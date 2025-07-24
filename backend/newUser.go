package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"golang.org/x/crypto/bcrypt"
)

func newUser(w http.ResponseWriter, r *http.Request) {
	userData, _ := fileSystem.GetUserData()
	if userData != "" {
		Coms.ExternalPostRespondCode(http.StatusForbidden, w)
		return
	}

	newUserPassword, err := Coms.ExternalPostReceived[string](r)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	userPasswordHash, err := createPasswordHash(*newUserPassword)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	err = fileSystem.SetUserData(string(userPasswordHash))
	if err != nil {
		Coms.PrintErrStr("Could not set user data in file system: ", err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
}

func createPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10) // 72 bytes max
}
