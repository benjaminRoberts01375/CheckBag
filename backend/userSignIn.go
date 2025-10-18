package main

import (
	"net/http"
	"time"

	"github.com/benjaminRoberts01375/CheckBag/backend/jwt"
	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
	"github.com/benjaminRoberts01375/CheckBag/backend/models"
	"golang.org/x/crypto/bcrypt"
)

func userSignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawPassword, err := requestReceived[string](r)
		if err != nil {
			Printing.PrintErrStr("Could not get password from request: ", err.Error())
			requestRespondCode(w, http.StatusBadRequest)
			return
		}
		passwordHash, err := fileSystem.GetUserData()
		if err != nil {
			Printing.PrintErrStr("Could not get user data from file system: ", err.Error())
			requestRespondCode(w, http.StatusBadRequest)
			return
		}
		// Compare the password with the hash
		if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(*rawPassword)); err != nil {
			Printing.PrintErrStr("Passwords do not match: ", err.Error())
			requestRespondCode(w, http.StatusBadRequest) // Intentionally obscure the error to prevent username guessing
			return
		}
		setJWTCookie(w)
		requestRespondCode(w, http.StatusOK)
	}
}

func setJWTCookie(w http.ResponseWriter) {
	// If the password is correct, generate a JWT
	token, err := jwt.GenerateJWT(jwt.LoginDuration)
	if err != nil {
		Printing.PrintErrStr("Could not generate JWT: ", err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}
	go cache.setUserSignIn(token)

	if models.ModelsConfig.DevMode {
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
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(jwt.LoginDuration),
			Path:     "/",
		})
	}
}
