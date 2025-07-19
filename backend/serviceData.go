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
	serviceData := []ServiceData{}

	// Handle time step requests
	// If there's a request for time steps, we're going to return a list of all services
	if len(queryParams["time-step"]) > 0 {
		serviceData = make([]ServiceData, len(serviceLinks))
		for i, service := range serviceLinks {
			serviceData[i] = ServiceData{ServiceLink: service, Hour: map[int]Analytic{}, Day: map[int]Analytic{}, Month: map[int]Analytic{}}
		}
	}

	for _, timeStepQuery := range queryParams["time-step"] {
		var timeStep AnalyticsTimeStep
		switch strings.ToLower(timeStepQuery) {
		case "day":
			timeStep = cacheAnalyticsDay
		case "month":
			timeStep = cacheAnalyticsMonth
		default:
			timeStep = cacheAnalyticsHour
		}
		cache.getAnalyticsTimeStep(timeStep, &serviceData)
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
