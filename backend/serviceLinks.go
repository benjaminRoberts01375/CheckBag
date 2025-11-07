package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
)

type ServiceLinks []ServiceLink

type ServiceLink struct {
	OutgoingAddress   ServiceAddress `json:"outgoing_address"`
	IncomingAddresses []string       `json:"incoming_addresses"`
	Title             string         `json:"title"`
	ID                string         `json:"id"`
}

type ServiceAddress struct {
	Protocol string `json:"protocol"`
	Domain   string `json:"domain"`
	Port     int    `json:"port"`
}

func (address ServiceAddress) String() string {
	return fmt.Sprintf("%s://%s:%d", address.Protocol, address.Domain, address.Port)
}

func (serviceLinks *ServiceLinks) Setup(db AdvancedDB) {
	ctx := context.Background()

	// Try to get services from database
	dbServices, err := db.getServiceLinks(ctx)
	if err != nil {
		// If database doesn't have services, try to migrate from filesystem
		Printing.PrintErrStr("Could not get services from database, attempting filesystem migration: " + err.Error())
		diskServices, fsErr := db.getServiceLinks(context.Background())
		if fsErr != nil {
			Printing.PrintErrStr("Could not get services from filesystem: " + fsErr.Error())
			return
		}

		// Migrate to database
		if len(diskServices) > 0 {
			migrateErr := db.setServiceLinks(ctx, diskServices)
			if migrateErr != nil {
				Printing.PrintErrStr("Could not migrate services to database: " + migrateErr.Error())
			} else {
				Printing.Println("Successfully migrated services from filesystem to database")
			}
		}
		*serviceLinks = diskServices
	} else {
		*serviceLinks = dbServices
	}

	Printing.Println("Loaded services: ", serviceLinks)
}

func (serviceLinks *ServiceLinks) String() string {
	var retVal string
	for _, serviceLink := range *serviceLinks {
		for _, incomingAddress := range serviceLink.IncomingAddresses {
			retVal += incomingAddress + " → " + serviceLink.OutgoingAddress.String() + "\n"
		}
	}
	return retVal
}

func servicesSet(serviceLinks *ServiceLinks, db AdvancedDB, jwt JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check JWT
		newServiceLinks, err := formatUserRequest[ServiceLinks](r, jwt)
		if err != nil {
			Printing.PrintErrStr("Could not add service: " + err.Error())
			requestRespondCode(w, http.StatusForbidden)
			return
		}

		// Delete service links that are not in new service links
		*serviceLinks = slices.DeleteFunc(*serviceLinks, func(existingService ServiceLink) bool {
			delVal := !slices.ContainsFunc(*newServiceLinks, func(newService ServiceLink) bool {
				return existingService.ID == newService.ID
			})
			if delVal {
				err := db.deleteService(r.Context(), existingService)
				if err != nil {
					Printing.PrintErrStr("Could not delete service from cache: " + err.Error())
				}
			}
			return delVal
		})

		// Update or add services to serviceLinks
		for _, newService := range *newServiceLinks {
			existingServiceI := slices.IndexFunc(*serviceLinks, func(existingService ServiceLink) bool {
				return existingService.ID == newService.ID
			})
			if existingServiceI == -1 { // Add service
				newService.ID = generateRandomString(15)
				*serviceLinks = append(*serviceLinks, newService)
				continue
			}
			// Update service - Don't update ID
			(*serviceLinks)[existingServiceI].IncomingAddresses = newService.IncomingAddresses
			(*serviceLinks)[existingServiceI].Title = newService.Title
			(*serviceLinks)[existingServiceI].OutgoingAddress = newService.OutgoingAddress
		}

		err = db.setServiceLinks(r.Context(), *serviceLinks)
		if err != nil {
			Printing.PrintErrStr("Could not set services in database: " + err.Error())
			requestRespondCode(w, http.StatusInternalServerError)
			return
		}
		Printing.Println("Updated service links: ", serviceLinks)
		requestRespond(w, serviceLinks)
	}
}

// Search for a service by incoming URL
func (services *ServiceLinks) GetServiceFromIncomingURL(service string) (*ServiceLink, error) {
	for _, serviceLink := range *services { // Check all services
		if slices.Contains(serviceLink.IncomingAddresses, service) {
			return &serviceLink, nil
		}
	}
	return nil, errors.New("no service found")
}

// Search for a service by outgoing URL
func (services *ServiceLinks) GetServiceFromOutgoingURL(service string) (*ServiceLink, error) {
	for _, serviceLink := range *services {
		if serviceLink.OutgoingAddress.String() == service {
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
