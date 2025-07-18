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
	// Services setup
	serviceLinks.Setup()
	// Analytics setup
	go startAnalyticsAdvance()
	// Setup endpoints
	setupEndpoints()
	if models.Config.DevMode {
		Coms.Println("Running in dev mode")
	}
	Coms.Println("Listening on port " + Coms.GetLaunchPort())
	http.ListenAndServe(Coms.GetLaunchPort(), nil)
}

func setupEndpoints() {
	http.HandleFunc("GET /api/user-exists", userExists)
	http.HandleFunc("POST /api/user-sign-up", newUser)
	http.HandleFunc("POST /api/user-sign-in", userSignIn)
	http.HandleFunc("POST /api/user-sign-in-jwt", userJWTSignIn)
	http.HandleFunc("POST /api/user-logout", userLogout)
	http.HandleFunc("POST /api/user-reset-password", userResetPassword)
	http.HandleFunc("POST /api/services-set", servicesSet)
	http.HandleFunc("GET /api/service-data/{service}", getServiceData)
	http.HandleFunc("/api/service/{path...}", requestForwarding)
	http.HandleFunc("/", notFound)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	Coms.PrintErrStr("404 - Not found - Sadge")
	w.Write([]byte("Womp Womp - 404"))
}
