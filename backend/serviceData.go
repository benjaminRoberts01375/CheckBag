package main

import (
	"net/http"
	"strings"
	"time"

	Printing "github.com/benjaminRoberts01375/Web-Tech-Stack/logging"
)

type ServiceData struct {
	Hour  map[time.Time]Analytic `json:"hour"`
	Day   map[time.Time]Analytic `json:"day"`
	Month map[time.Time]Analytic `json:"month"`
	Year  map[time.Time]Analytic `json:"year"`
	ServiceLink
}

type Analytic struct {
	Quantity     int            `json:"quantity"`
	Country      map[string]int `json:"country"`
	IP           map[string]int `json:"ip"`
	Resource     map[string]int `json:"resource"`
	ResponseCode map[int]int    `json:"response_code"`
}

func getServiceData(w http.ResponseWriter, r *http.Request) {
	_, _, err := checkUserRequest[any](r)
	if err != nil {
		Printing.PrintErrStr("Could not verify user for analytic data: " + err.Error())
		requestRespondCode(w, http.StatusForbidden)
		return
	}
	queryParams := r.URL.Query()
	serviceData := make([]ServiceData, len(serviceLinks))

	// Create a list of all services
	for i, service := range serviceLinks {
		serviceData[i] = ServiceData{ServiceLink: service, Hour: map[time.Time]Analytic{}, Day: map[time.Time]Analytic{}, Month: map[time.Time]Analytic{}, Year: map[time.Time]Analytic{}}
	}

	// Handle time step requests
	for _, timeScaleQuery := range queryParams["time-step"] {
		timeScaleQuery = strings.ToLower(timeScaleQuery)
		switch timeScaleQuery {
		case "month":
			for i, service := range serviceData {
				serviceData[i].Month = cache.getAnalyticsService(service, cacheAnalyticsDay)
			}
		case "year":
			for i, service := range serviceData {
				serviceData[i].Year = cache.getAnalyticsService(service, cacheAnalyticsMonth)
			}
		case "day":
			for i, service := range serviceData {
				serviceData[i].Day = cache.getAnalyticsService(service, cacheAnalyticsHour)
			}
		default: // hour
			for i, service := range serviceData {
				serviceData[i].Hour = cache.getAnalyticsService(service, cacheAnalyticsMinute)
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
		serviceData[requestedServiceIndex].Day = cache.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsHour)
		serviceData[requestedServiceIndex].Month = cache.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsDay)
		serviceData[requestedServiceIndex].Year = cache.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsMonth)
		serviceData[requestedServiceIndex].Hour = cache.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsMinute)
	}

	requestRespond(w, serviceData)
}
