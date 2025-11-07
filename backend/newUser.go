package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func newUser(db AdvancedDB, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, _ := db.GetUserPasswordHash(r.Context())
		if userData != "" {
			requestRespondCode(w, http.StatusForbidden)
			return
		}

		newUserPassword, err := requestReceived[string](r)
		if err != nil {
			requestRespondCode(w, http.StatusBadRequest)
			return
		}
		userPasswordHash, err := createPasswordHash(*newUserPassword)
		if err != nil {
			requestRespondCode(w, http.StatusBadRequest)
			return
		}
		db.SetUserPasswordHash(r.Context(), string(userPasswordHash))
		jwt.setJWT(w)
		requestRespondCode(w, http.StatusOK)
	}
}
func createPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 10) // 72 bytes max
}
