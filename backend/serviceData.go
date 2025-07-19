package main

import (
	"net/http"
	"strings"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

type ServiceData struct {
	Hour  map[int]Analytic `json:"hour"`
	Day   map[int]Analytic `json:"day"`
	Month map[int]Analytic `json:"month"`
	ServiceLink
}

type Analytic struct {
	Quantity int            `json:"quantity"`
	Country  map[string]int `json:"country"`
	IP       map[string]int `json:"ip"`
	Resource map[string]int `json:"resource"`
}

func getServiceData(w http.ResponseWriter, r *http.Request) {
	_, _, err := checkUserRequest[any](r)
	if err != nil {
		Coms.PrintErrStr("Could not verify user for analytic data: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	queryParams := r.URL.Query()
	serviceData := make([]ServiceData, len(serviceLinks))

	// Create a list of all services
	for i, service := range serviceLinks {
		serviceData[i] = ServiceData{ServiceLink: service, Hour: map[int]Analytic{}, Day: map[int]Analytic{}, Month: map[int]Analytic{}}
	}

	// Handle time step requests
	for _, timeStepQuery := range queryParams["time-step"] {
		timeStepQuery = strings.ToLower(timeStepQuery)
		switch timeStepQuery {
		case "day":
			for i, service := range serviceData {
				serviceData[i].Day = cache.getAnalyticsService(service, cacheAnalyticsDay)
			}
		case "month":
			for i, service := range serviceData {
				serviceData[i].Month = cache.getAnalyticsService(service, cacheAnalyticsMonth)
			}
		default: // Hour
			for i, service := range serviceData {
				serviceData[i].Hour = cache.getAnalyticsService(service, cacheAnalyticsHour)
			}
		}
	}

	// Handle service requests
	for _, requestedServiceString := range queryParams["service"] {
		// Find the service index of the requested service from the list of all service data
		requestedServiceLink, _ := serviceLinks.GetService(requestedServiceString)
		var requestedServiceIndex int = -1
		for i, service := range serviceData {
			if service.ServiceLink.ID == requestedServiceLink.ID {
				requestedServiceIndex = i
				break
			}
		}
		// Check if the requested service was found
		if requestedServiceIndex == -1 {
			continue
		}
		// Get the requested service's analytics
		serviceData[requestedServiceIndex].Hour = cache.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsHour)
		serviceData[requestedServiceIndex].Day = cache.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsDay)
		serviceData[requestedServiceIndex].Month = cache.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsMonth)
	}

	Coms.ExternalPostRespond(serviceData, w)
}
