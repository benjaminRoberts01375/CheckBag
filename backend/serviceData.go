package main

import (
	"net/http"

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
	requestedService := r.PathValue("services")
	if requestedService == "" {
		Coms.ExternalPostRespond(getAnalyticDashboard(), w)
		return
	} else {
		Coms.ExternalPostRespond(getAnalyticData(), w)
	}
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
