package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/models"
)

var cache CacheClient[*CacheLayer]

func main() {
	// Config setup
	models.Setup()
	// Coms setup
	Coms.ReadConfig()
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
	http.HandleFunc("POST /api/user-sign-in", newUserSignIn)
	http.HandleFunc("POST /api/user-sign-in-jwt", userJWTSignIn)
	http.HandleFunc("POST /api/user-logout", userLogout)
	http.HandleFunc("POST /api/user-reset-password", userResetPassword)
	http.HandleFunc("POST /api/user-get-data", userRequestData)
}
