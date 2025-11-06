package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
)

func main() {
	var fileSystem = FileSystem{
		BasePath: "/CheckBag",
		UserData: "userdata.txt",
		Services: "services.json",
	}
	fileSystem.Setup()

	var serviceLinks = ServiceLinks{}

	// Coms setup
	Printing.ReadConfig()
	// Cache setup
	db := SetupDB()
	// Services setup
	serviceLinks.Setup(fileSystem, db)
	// JWT Setup
	jwt := loadJWTSecret(db)
	// Setup endpoints
	setupEndpoints(fileSystem, &serviceLinks, db, jwt, strings.ToLower(os.Getenv("DEV_MODE")) == "true")
	Printing.Println("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}

func setupEndpoints(fileSystem FileSystem, serviceLinks *ServiceLinks, db AdvancedDB, jwt JWTService, devMode bool) {
	http.HandleFunc("GET /api/user-exists", userExists(fileSystem))                           // Check if the user already exists
	http.HandleFunc("POST /api/user-sign-up", newUser(fileSystem, jwt))                       // Sign up with username and password
	http.HandleFunc("POST /api/user-sign-in", userSignIn(fileSystem, jwt))                    // Sign in with username and password
	http.HandleFunc("POST /api/user-sign-in-jwt", userJWTSignIn(jwt))                         // Sign in with JWT
	http.HandleFunc("POST /api/services-set", servicesSet(fileSystem, serviceLinks, db, jwt)) // Setting/replacing all services
	http.HandleFunc("GET /api/service-data", getServiceData(serviceLinks, db, jwt))           // Getting analytics
	http.HandleFunc("/api/service/{path...}", requestForwarding(serviceLinks, db))            // Proxying requests
	http.HandleFunc("GET /api/api-keys", APIGet(db, jwt))                                     // Getting API keys
	http.HandleFunc("POST /api/api-keys", APISet(db, jwt))                                    // Setting API keys

	http.HandleFunc("/", spaHandler(devMode)) // Serve the frontend
}

// spaHandler serves the React SPA and handles client-side routing
func spaHandler(devMode bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Skip API routes - something's wrong :(
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}

		// Production mode: serve static files
		// Development mode: proxy to Vite dev server
		if devMode {
			proxyToVite(w, r)
			return
		}

		// Construct the path to the requested file
		path := filepath.Join("./static", r.URL.Path)

		// Check if the file exists
		if _, err := os.Stat(path); os.IsNotExist(err) {
			// File doesn't exist, serve index.html for React routing
			http.ServeFile(w, r, "./static/index.html")
			return
		}

		// File exists, serve it normally
		http.FileServer(http.Dir("./static")).ServeHTTP(w, r)
	}
}

// A quick and dirty proxy for the Vite dev server
func proxyToVite(w http.ResponseWriter, r *http.Request) {
	// Parse Vite dev server URL
	viteURL, err := url.Parse("http://frontend:5173")
	if err != nil {
		Printing.Printf("Failed to parse Vite URL: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create reverse proxy
	proxy := httputil.NewSingleHostReverseProxy(viteURL)

	// Add error handling for when Vite isn't ready
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		http.Error(w, "Frontend development server not ready", http.StatusBadGateway)
		Printing.PrintErrStr("Proxy error to Vite dev server:", err.Error())
	}

	// Modify the request headers for proper proxying
	r.URL.Host = viteURL.Host
	r.URL.Scheme = viteURL.Scheme
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
	r.Host = viteURL.Host

	// Proxy the request
	proxy.ServeHTTP(w, r)
}
