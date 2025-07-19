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

	for _, timeStepQuery := range queryParams["time-step"] { // Handle time step requests
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
		default:
			for i, service := range serviceData {
				serviceData[i].Hour = cache.getAnalyticsService(service, cacheAnalyticsHour)
			}
		}
	}

	Coms.ExternalPostRespond(serviceData, w)
}

func getAnalyticDashboard() []ServiceData {
	var serviceData []ServiceData = make([]ServiceData, len(serviceLinks))
	for i, service := range serviceLinks {
		serviceData[i] = ServiceData{
			ServiceLink: service,
		}
	}
	return serviceData
}

func getAnalyticData() ServiceData {
	return ServiceData{}
}
