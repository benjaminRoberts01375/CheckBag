package main

import (
	"errors"
	"net/http"
	"slices"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
	"github.com/benjaminRoberts01375/CheckBag/backend/models"
)

type ServiceLinks []ServiceLink

type ServiceLink struct {
	OutgoingAddress   string   `json:"outgoing_address"`
	IncomingAddresses []string `json:"incoming_address"`
	Title             string   `json:"title"`
	ID                string   `json:"id"`
}

func (serviceLinks *ServiceLinks) Setup() {
	diskServices, err := fileSystem.GetServices()
	if err != nil {
		Printing.PrintErrStr("Could not get services: " + err.Error())
		return
	}
	*serviceLinks = diskServices
	Printing.Println("Loaded services: ", serviceLinks)
}

func (serviceLinks *ServiceLinks) String() string {
	var retVal string
	for _, serviceLink := range *serviceLinks {
		retVal += serviceLink.IncomingAddresses[0] + " → " + serviceLink.OutgoingAddress + "\n"
	}
	return retVal
}

var serviceLinks = ServiceLinks{}

func servicesSet(w http.ResponseWriter, r *http.Request) {
	// Check JWT
	_, newServiceLinks, err := checkUserRequest[ServiceLinks](r)
	if err != nil {
		Printing.PrintErrStr("Could not add service: " + err.Error())
		requestRespondCode(w, http.StatusForbidden)
		return
	}

	// Delete service links that are not in new service links
	serviceLinks = slices.DeleteFunc(serviceLinks, func(existingService ServiceLink) bool {
		delVal := !slices.ContainsFunc(*newServiceLinks, func(newService ServiceLink) bool {
			return existingService.ID == newService.ID
		})
		if delVal {
			err := cache.deleteService(existingService)
			if err != nil {
				Printing.PrintErrStr("Could not delete service from cache: " + err.Error())
			}
		}
		return delVal
	})

	// Update or add services to serviceLinks
	for _, newService := range *newServiceLinks {
		existingServiceI := slices.IndexFunc(serviceLinks, func(existingService ServiceLink) bool {
			return existingService.ID == newService.ID
		})
		if existingServiceI == -1 { // Add service
			newService.ID = generateRandomString(models.ModelsConfig.ServiceIDLength)
			serviceLinks = append(serviceLinks, newService)
			continue
		}
		// Update service - Don't update ID
		serviceLinks[existingServiceI].IncomingAddresses = newService.IncomingAddresses
		serviceLinks[existingServiceI].Title = newService.Title
		serviceLinks[existingServiceI].OutgoingAddress = newService.OutgoingAddress
	}

	err = fileSystem.SetServices(serviceLinks)
	if err != nil {
		Printing.PrintErrStr("Could not set services in file system: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}
	Printing.Println("Updated service links")
	requestRespond(w, serviceLinks)
}

// Search for a service by incoming URL
func (services *ServiceLinks) GetServiceFromIncomingURL(service string) (*ServiceLink, error) {
	for _, serviceLink := range *services {
		if slices.Contains(serviceLink.IncomingAddresses, service) {
			return &serviceLink, nil
		}
	}
	return nil, errors.New("no service found")
}

// Search for a service by outgoing URL
func (services *ServiceLinks) GetServiceFromOutgoingURL(service string) (*ServiceLink, error) {
	for _, serviceLink := range *services {
		if serviceLink.OutgoingAddress == service {
			return &serviceLink, nil
		}
	}
	return nil, errors.New("no service found")
}

// Search for a service by ID
func (services *ServiceLinks) GetServiceByID(serviceID string) (*ServiceLink, error) {
	for _, serviceLink := range *services {
		if serviceLink.ID == serviceID {
			return &serviceLink, nil
		}
	}
	return nil, errors.New("no service found")
}
