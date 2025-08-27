package main

import (
	"net/http"
	"strings"
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
		} else if ip != "" && strings.HasSuffix(name, "ip") { // Get the client's IP address
			ip = values[0]
		} else if name == "x-forwarded-for" { // Get the client's IP address
			ip = values[0]
		}
	}
	cache.incrementAnalytics(serviceID, resource, country, ip, responseCode)
}
