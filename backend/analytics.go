package main

import (
	"net/http"
	"strings"
	"time"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

func analytics(r *http.Request, responseCode int) {
	service, err := serviceLinks.GetServiceFromExternalURL(r.Host)
	var serviceID string
	if err != nil {
		serviceID = "Unknown"
	} else {
		serviceID = service.ID
	}
	resource := r.PathValue("path")

	// Search for country and ip in headers
	var country string
	var ip string
	for name, values := range r.Header {
		name = strings.ToLower(name)
		// Cloudflare headers have both "country" and "ip" in the header name for country
		if strings.Contains(name, "country") {
			country = values[0]
		} else if strings.Contains(name, "ip") {
			ip = values[0]
		}
	}
	cache.incrementAnalytics(serviceID, resource, country, ip, responseCode)
}

func startAnalyticsAdvance() {
	Coms.Println("Starting analytics advance")
	triggerChan := make(chan AnalyticsTimeStep, len(cacheAnalyticsTime))

	for _, timeStep := range cacheAnalyticsTime {
		go func(timeStep AnalyticsTimeStep) {
			remainingTime := timeStep.timeToNextStep() // Prevent drift
			ticker := time.NewTicker(remainingTime)
			defer ticker.Stop()

			for range ticker.C {
				ticker.Reset(timeStep.timeToNextStep()) // Prevent drift
				triggerChan <- timeStep
			}
		}(timeStep)
	}

	for source := range triggerChan {
		cache.advanceAnalytics(source, serviceLinks)
	}
	Coms.Println("Finished analytics advance")
}
