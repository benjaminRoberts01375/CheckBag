package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/jwt"
)

func userLogout(w http.ResponseWriter, r *http.Request) {
	_, _, _, err := checkUserRequest[any](r)
	if err != nil {
		Coms.PrintErrStr("Could not log out user: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	cookie, err := r.Cookie(jwt.CookieName)
	if err != nil {
		Coms.PrintErrStr("Could not get user JWT cookie: " + err.Error())
		return
	}

	err = cache.deleteUserSignIn(cookie.Value)
	if err != nil {
		Coms.PrintErrStr("Could not delete user JWT: " + err.Error())
		return
	}
	Coms.ExternalPostRespondCode(http.StatusOK, w)
}
