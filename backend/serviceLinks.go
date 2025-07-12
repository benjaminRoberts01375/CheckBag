package main

import (
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

type ServiceLinks []ServiceLink

type ServiceLink struct {
	InternalAddress string   `json:"internal_address"`
	ExternalAddress []string `json:"external_address"`
	Title           string   `json:"title"`
}

func (serviceLinks *ServiceLinks) Setup() {
	diskServices, err := fileSystem.GetServices()
	if err != nil {
		Coms.PrintErrStr("Could not get services: " + err.Error())
		return
	}
	*serviceLinks = diskServices
}

var serviceLinks = ServiceLinks{}

func servicesSet(w http.ResponseWriter, r *http.Request) {
	// Check JWT
	_, newServiceLinks, err := checkUserRequest[ServiceLinks](r)
	if err != nil {
		Coms.PrintErrStr("Could not add service: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	err = fileSystem.SetServices(*newServiceLinks)
	if err != nil {
		Coms.PrintErrStr("Could not set services in file system: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusInternalServerError, w)
		return
	}
	serviceLinks = *newServiceLinks
}
