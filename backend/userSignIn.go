package main

import (
	"net/http"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
	"golang.org/x/crypto/bcrypt"
)

func userSignIn(db AdvancedDB, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawPassword, err := requestReceived[string](r)
		if err != nil {
			Printing.PrintErrStr("Could not get password from request: ", err.Error())
			requestRespondCode(w, http.StatusBadRequest)
			return
		}
		passwordHash, err := db.GetUserPasswordHash(r.Context())
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
		jwt.setJWT(w)
		requestRespondCode(w, http.StatusOK)
	}
}
