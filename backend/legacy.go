package main

import (
	"strconv"
	"strings"
)

// *** Legacy data types and conversions ***

// Service Link
type ServiceLinkV1 struct {
	InternalAddress string   `json:"internal_address"`
	ExternalAddress []string `json:"external_address"`
	Title           string   `json:"title"`
	ID              string   `json:"id"`
}

func (serviceLink ServiceLinkV1) Migrate() ServiceLink {
	internalAddressParts := strings.Split(serviceLink.InternalAddress, ":")
	if len(internalAddressParts) != 2 {
		return ServiceLink{}
	} else if len(serviceLink.ExternalAddress) != 1 {
		return ServiceLink{}
	}
	internalPort, err := strconv.Atoi(internalAddressParts[1])
	if err != nil {
		return ServiceLink{}
	}

	return ServiceLink{
		IncomingAddresses: serviceLink.ExternalAddress,
		OutgoingAddress: ServiceAddress{ // Always uses http on some port
			Protocol: "http",
			Domain:   internalAddressParts[0],
			Port:     internalPort,
		},
		Title: serviceLink.Title,
		ID:    serviceLink.ID,
	}
}

// Service Links
type ServiceLinksV1 []ServiceLinkV1

func (serviceLinks ServiceLinksV1) Migrate() ServiceLinks {
	var migratedServiceLinks ServiceLinks
	for _, serviceLink := range serviceLinks {
		migratedServiceLinks = append(migratedServiceLinks, serviceLink.Migrate())
	}
	return migratedServiceLinks
}
