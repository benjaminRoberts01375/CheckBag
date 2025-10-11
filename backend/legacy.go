package main

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
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

// Merge all v2 service links
//
// mergeAllServiceLinks merges all v2 service links that have the same exact
// outgoing address (protocol, domain, and port). It also merges the analytics
// stored in Valkey. This is useful for cleaning up duplicate services that may
// have been created in older versions of CheckBag.
func (serviceLinks *ServiceLinks) mergeAllServiceLinks() error {
	// Group services by outgoing address
	serviceGroups := make(map[string][]int) // map[outgoingAddress][]serviceIndex

	for i, service := range *serviceLinks {
		key := service.OutgoingAddress.String()
		serviceGroups[key] = append(serviceGroups[key], i)
	}

	// Track services to remove
	servicesToRemove := make([]int, 0)

	// Process each group that has duplicates
	for _, indices := range serviceGroups {
		if len(indices) <= 1 {
			continue // No merge needed
		}

		// Keep the first service, merge others into it
		keepIndex := indices[0]
		keepService := &(*serviceLinks)[keepIndex]

		Printing.Printf("Merging %d duplicate services with outgoing address %s",
			len(indices), keepService.OutgoingAddress.String())

		// Merge incoming addresses and analytics from other services
		for _, removeIndex := range indices[1:] {
			removeService := (*serviceLinks)[removeIndex]

			// Merge incoming addresses (avoid duplicates)
			for _, incomingAddr := range removeService.IncomingAddresses {
				if !slices.Contains(keepService.IncomingAddresses, incomingAddr) {
					keepService.IncomingAddresses = append(keepService.IncomingAddresses, incomingAddr)
					Printing.Printf("Added incoming address %s to service %s(%s)", incomingAddr, keepService.Title, keepService.ID)
				}
			}

			// Merge analytics in Valkey
			Printing.Printf("Merging analytics from service %s(%s) to %s(%s)", removeService.Title, removeService.ID, keepService.Title, keepService.ID)
			if err := mergeServiceAnalytics(keepService.ID, removeService.ID); err != nil {
				Printing.PrintErrStr("Failed to merge analytics:", err.Error())
				return fmt.Errorf("failed to merge analytics from %s to %s: %w",
					removeService.ID, keepService.ID, err)
			}

			// Mark for removal
			servicesToRemove = append(servicesToRemove, removeIndex)
		}

		Printing.Printf("Successfully merged into service %s (%s)",
			keepService.Title, keepService.ID)
	}

	if len(servicesToRemove) == 0 {
		Printing.Println("No duplicate services found to merge")
		return nil
	}

	// Remove merged services (in reverse order to maintain indices)
	slices.Sort(servicesToRemove)
	slices.Reverse(servicesToRemove)
	for _, index := range servicesToRemove {
		removedService := (*serviceLinks)[index]
		Printing.Printf("Removing service %s (%s)", removedService.Title, removedService.ID)
		*serviceLinks = slices.Delete(*serviceLinks, index, index+1)
	}

	Printing.Printf("Merge complete: removed %d duplicate services", len(servicesToRemove))
	return nil
}

// mergeServiceAnalytics merges all analytics data from removeServiceID into keepServiceID
func mergeServiceAnalytics(keepServiceID, removeServiceID string) error {
	if keepServiceID == removeServiceID {
		return nil // Nothing to merge
	}

	mergedKeys := 0

	// For each time step (minute, hour, day, month)
	for _, timeStep := range cacheAnalyticsTime {
		quantity := strconv.Itoa(timeStep.maximumUnits)

		// For each time period in the time step
		for timePeriod := range timeStep.maximumUnits {
			timestamp := timeStep.timeStr(-timePeriod)
			expiration := timeStep.time(timeStep.maximumUnits - timePeriod)

			// Merge quantity (simple counter)
			removeQuantityKey := fmt.Sprintf("Analytics:%s:%s:%s:quantity", removeServiceID, quantity, timestamp)
			removeQuantityStr, err := cache.raw.Get(removeQuantityKey)
			if err == nil {
				removeQuantityInt, err := strconv.Atoi(removeQuantityStr)
				if err == nil && removeQuantityInt > 0 {
					keepQuantityKey := fmt.Sprintf("Analytics:%s:%s:%s:quantity", keepServiceID, quantity, timestamp)
					err = cache.raw.IncrementKeyBy(keepQuantityKey, removeQuantityInt, expiration)
					if err != nil {
						return fmt.Errorf("failed to merge quantity: %w", err)
					}
					mergedKeys++
				}
				// Delete the old key
				cache.raw.Delete(removeQuantityKey)
			}

			// Merge hash fields (country, ip, resource, response_code)
			hashTypes := []string{"country", "ip", "resource", "response_code"}
			for _, hashType := range hashTypes {
				removeKey := fmt.Sprintf("Analytics:%s:%s:%s:%s", removeServiceID, quantity, timestamp, hashType)
				removeHash, err := cache.raw.GetHash(removeKey)
				if err == nil && len(removeHash) > 0 {
					keepKey := fmt.Sprintf("Analytics:%s:%s:%s:%s", keepServiceID, quantity, timestamp, hashType)

					// Increment each field in the keep service's hash
					for field, valueStr := range removeHash {
						valueInt, err := strconv.Atoi(valueStr)
						if err == nil && valueInt > 0 {
							err = cache.raw.IncrementHashField(keepKey, field, valueInt, expiration)
							if err != nil {
								return fmt.Errorf("failed to merge %s field %s: %w", hashType, field, err)
							}
						}
					}

					mergedKeys++
					// Delete the old hash
					cache.raw.DeleteHash(removeKey)
				}
			}
		}
	}

	if mergedKeys > 0 {
		Printing.Printf("Merged %d analytics keys", mergedKeys)
	}

	return nil
}

// Service Links
type ServiceLinksV1 []ServiceLinkV1

func (serviceLinks ServiceLinksV1) Migrate() ServiceLinks {
	var migratedServiceLinks ServiceLinks
	for _, serviceLink := range serviceLinks {
		migratedServiceLinks = append(migratedServiceLinks, serviceLink.Migrate())
	}
	err := migratedServiceLinks.mergeAllServiceLinks()
	if err != nil {
		Printing.PrintErrStr("Failed to merge service links:", err.Error())
	}
	return migratedServiceLinks
}
