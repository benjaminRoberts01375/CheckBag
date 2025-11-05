package main

import (
	"net/http"
)

func userJWTSignIn(db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, err := checkUserRequest[any](r, db)
		if err != nil {
			requestRespondCode(w, http.StatusBadRequest)
			return
		}
		requestRespondCode(w, http.StatusOK)
	}
}
