package main

import (
	"net/http"
)

func userJWTSignIn(jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := jwt.ReadAndValidateJWT(r)
		if err != nil {
			requestRespondCode(w, http.StatusBadRequest)
			return
		}
		requestRespondCode(w, http.StatusOK)
	}
}
