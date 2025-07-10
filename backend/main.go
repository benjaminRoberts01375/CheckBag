package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/models"
)

var database *DBClient
var cache CacheClient[*CacheLayer]

func main() {
	// Config setup
	models.Setup()
	// Coms setup
	Coms.ReadConfig()
	// Database setup
	database = &DBClient{}
	database.raw = DBLayer{}
	database.raw.Setup()
	defer database.raw.Close()
	// Cache setup
	cache.raw = &CacheLayer{}
	cache.raw.Setup()
	defer cache.raw.Close()
	// Setup endpoints
	setupEndpoints()

	http.ListenAndServe(Coms.GetLaunchPort(), nil)
}

func setupEndpoints() {
	http.HandleFunc("POST /api/user-sign-up", newUser)
	http.HandleFunc("/api/user-sign-up-confirmation/{token}", newUserConfirmation)
	http.HandleFunc("POST /api/user-sign-in", newUserSignIn)
	http.HandleFunc("POST /api/user-sign-in-jwt", userJWTSignIn)
	http.HandleFunc("POST /api/user-logout", userLogout)
	http.HandleFunc("POST /api/user-reset-password", userResetPassword)
	http.HandleFunc("POST /api/user-reset-email-request", userResetEmailRequest)
	http.HandleFunc("/api/user-reset-email-confirmation/{token}", userResetEmailConfirmation)
	http.HandleFunc("POST /api/user-forgot-password-request", userForgotPasswordRequest)
	http.HandleFunc("POST /api/user-forgot-password-check-valid/{token}", userForgotPasswordCheckValid)
	http.HandleFunc("POST /api/user-forgot-password-confirmation/{token}", userForgotPasswordConfirmation)
	http.HandleFunc("POST /api/user-get-data", userRequestData)
}
