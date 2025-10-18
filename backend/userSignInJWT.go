package main

import (
	"net/http"
)

func userJWTSignIn(cache CacheClient[*CacheLayer]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, err := checkUserRequest[any](r, cache)
		if err != nil {
			requestRespondCode(w, http.StatusBadRequest)
			return
		}
		requestRespondCode(w, http.StatusOK)
	}
}
