package main

import (
	"net/http"
	"strings"
	"time"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
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

func getServiceData(serviceLinks *ServiceLinks, db AdvancedDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		_, _, err := checkUserRequest[any](r, db)
		if err != nil {
			if !(len(queryParams["api-key"]) > 0 && db.apiKeyExists(queryParams["api-key"][0])) { // Check if API key is invalid
				Printing.PrintErrStr("Could not verify user or API key for analytic data: " + err.Error())
				requestRespondCode(w, http.StatusForbidden)
				return
			}
		}

		serviceData := make([]ServiceData, len(*serviceLinks))

		// Create a list of all services
		for i, service := range *serviceLinks {
			serviceData[i] = ServiceData{ServiceLink: service, Hour: map[time.Time]Analytic{}, Day: map[time.Time]Analytic{}, Month: map[time.Time]Analytic{}, Year: map[time.Time]Analytic{}}
		}

		// Handle time step requests
		for _, timeScaleQuery := range queryParams["time-step"] {
			timeScaleQuery = strings.ToLower(timeScaleQuery)
			switch timeScaleQuery {
			case "hour":
				for i, service := range serviceData {
					serviceData[i].Hour = db.getAnalyticsService(service, cacheAnalyticsMinute)
				}
			case "day":
				for i, service := range serviceData {
					serviceData[i].Day = db.getAnalyticsService(service, cacheAnalyticsHour)
				}
			case "month":
				for i, service := range serviceData {
					serviceData[i].Month = db.getAnalyticsService(service, cacheAnalyticsDay)
				}
			case "year":
				for i, service := range serviceData {
					serviceData[i].Year = db.getAnalyticsService(service, cacheAnalyticsMonth)
				}
			}
		}

		// Handle service requests
		for _, requestedServiceString := range queryParams["service"] {
			// Find the service index of the requested service from the list of all service data
			requestedServiceLink, _ := serviceLinks.GetServiceFromIncomingURL(requestedServiceString)
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
			serviceData[requestedServiceIndex].Day = db.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsHour)
			serviceData[requestedServiceIndex].Month = db.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsDay)
			serviceData[requestedServiceIndex].Year = db.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsMonth)
			serviceData[requestedServiceIndex].Hour = db.getAnalyticsService(serviceData[requestedServiceIndex], cacheAnalyticsMinute)
		}

		requestRespond(w, serviceData)
	}
}
