package main

import (
	"net/http"
)

func userJWTSignIn(w http.ResponseWriter, r *http.Request) {
	_, _, err := checkUserRequest[any](r)
	if err != nil {
		requestRespondCode(w, http.StatusBadRequest)
		return
	}
	requestRespondCode(w, http.StatusOK)
}
