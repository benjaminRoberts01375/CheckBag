package main

import (
	"net/http"
	"time"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/jwt"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/models"
	"golang.org/x/crypto/bcrypt"
)

func userSignIn(w http.ResponseWriter, r *http.Request) {
	rawPassword, err := Coms.ExternalPostReceived[string](r)
	if err != nil {
		Coms.PrintErrStr("Could not get password from request: ", err.Error())
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	passwordHash, err := fileSystem.GetUserData()
	if err != nil {
		Coms.PrintErrStr("Could not get user data from file system: ", err.Error())
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	// Compare the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(*rawPassword)); err != nil {
		Coms.PrintErrStr("Passwords do not match: ", err.Error())
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w) // Intentionally obscure the error to prevent username guessing
		return
	}

	// If the password is correct, generate a JWT
	token, err := jwt.GenerateJWT(jwt.LoginDuration)
	if err != nil {
		Coms.PrintErrStr("Could not generate JWT: ", err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	go cache.setUserSignIn(token)

	if models.Config.DevMode {
		http.SetCookie(w, &http.Cookie{
			Name:     jwt.CookieName,
			Value:    token,
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(jwt.LoginDuration),
			Path:     "/",
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:     jwt.CookieName,
			Value:    token,
			HttpOnly: false,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(jwt.LoginDuration),
			Path:     "/",
		})
	}
	Coms.ExternalPostRespondCode(http.StatusOK, w)
}
