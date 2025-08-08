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
		} else if strings.Contains(name, "ip") {
			ip = values[0]
		}
	}
	serviceUpdater.Update <- RequestAnalyticData{
		ServiceID:    serviceID,
		Resource:     resource,
		Country:      country,
		IP:           ip,
		ResponseCode: responseCode,
	}
	cache.incrementAnalytics(serviceID, resource, country, ip, responseCode)
}
