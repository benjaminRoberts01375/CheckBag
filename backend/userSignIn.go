package main

import (
	"net/http"
	"time"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/jwt"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/models"
	"golang.org/x/crypto/bcrypt"
)

type UserSignIn struct {
	Email    string `json:"username"`
	Password string `json:"password"`
}

func newUserSignIn(w http.ResponseWriter, r *http.Request) {
	userRequest, err := Coms.ExternalPostReceived[UserSignIn](r)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	dbPasswordHash, userID, err := database.GetUserPasswordAndID(userRequest.Email)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w)
		return
	}
	// Compare the password with the hash
	if err := bcrypt.CompareHashAndPassword(dbPasswordHash, []byte(userRequest.Password)); err != nil {
		Coms.ExternalPostRespondCode(http.StatusBadRequest, w) // Intentionally obscure the error to prevent username guessing
		return
	}

	// If the password is correct, generate a JWT
	token, err := jwt.GenerateJWT(userRequest.Email, jwt.LoginDuration)
	if err != nil {
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	go database.SetUserHasLoggedIn(userID)
	go cache.setUserSignIn(token, userID)

	if models.Config.DevMode {
		http.SetCookie(w, &http.Cookie{
			Name:     jwt.CookieName,
			Value:    token,
			HttpOnly: false,
			Secure:   false,
			SameSite: http.SameSiteStrictMode,
			Expires:  time.Now().Add(time.Hour * 24 * 30), // One month
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
