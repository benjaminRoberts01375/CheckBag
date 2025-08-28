package main

import (
	"net/http"

	"github.com/benjaminRoberts01375/CheckBag/backend/jwt"
	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
)

func userLogout(w http.ResponseWriter, r *http.Request) {
	_, _, err := checkUserRequest[any](r)
	if err != nil {
		Printing.PrintErrStr("Could not log out user: " + err.Error())
		requestRespondCode(w, http.StatusForbidden)
		return
	}
	cookie, err := r.Cookie(jwt.CookieName)
	if err != nil {
		Printing.PrintErrStr("Could not get user JWT cookie: " + err.Error())
		return
	}

	err = cache.deleteUserSignIn(cookie.Value)
	if err != nil {
		Printing.PrintErrStr("Could not delete user JWT: " + err.Error())
		return
	}
	requestRespondCode(w, http.StatusOK)
}
