package main

// *** Legacy data types and conversions ***

// Service Link
type ServiceLinkV1 struct {
	InternalAddress string   `json:"internal_address"`
	ExternalAddress []string `json:"external_address"`
	Title           string   `json:"title"`
	ID              string   `json:"id"`
}

func (serviceLink ServiceLinkV1) Migrate() ServiceLink {
	return ServiceLink{
		IncomingAddresses: serviceLink.ExternalAddress,
		OutgoingAddress:   serviceLink.InternalAddress,
		Title:             serviceLink.Title,
		ID:                serviceLink.ID,
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
