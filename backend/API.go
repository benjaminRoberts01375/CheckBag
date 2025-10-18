package main

import (
	"net/http"
	"slices"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
)

type APIKeyInfo struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	ID   string `json:"id"`
}

func APISet(cache CacheClient[*CacheLayer]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Ensure user is logged in
		_, newKeys, err := checkUserRequest[[]APIKeyInfo](r, cache)
		if err != nil {
			Printing.PrintErrStr("Could not create API: " + err.Error())
			requestRespondCode(w, http.StatusForbidden)
			return
		}

		// Get all existing keys and their data
		existingKeys, err := cache.getAPIKeyInfo()
		if err != nil {
			Printing.PrintErrStr("Could not get API keys: " + err.Error())
			requestRespondCode(w, http.StatusInternalServerError)
			return
		}

		// Remove all keys from cache that are not in newKeys
		for _, existingKey := range existingKeys {
			if !slices.ContainsFunc(*newKeys, func(newKey APIKeyInfo) bool {
				return existingKey.ID == newKey.ID // Check if the key already exists in newKeys
			}) {
				err := cache.removeAPIKey(existingKey.ID)
				if err != nil {
					Printing.PrintErrStr("Could not remove API key: " + err.Error())
					requestRespondCode(w, http.StatusInternalServerError)
					return
				}
			}
		}

		// Add all new keys to cache
		for i, newKey := range *newKeys {
			if !slices.ContainsFunc(existingKeys, func(existingKey APIKeyInfo) bool {
				return existingKey.ID == newKey.ID // Check if the key already exists in existingKeys
			}) {
				APIKey := generateRandomString(32)
				keyID := generateRandomString(32)
				(*newKeys)[i].Key = APIKey
				(*newKeys)[i].ID = keyID
				err = cache.addAPIKey(APIKey, keyID, newKey.Name)
				if err != nil {
					Printing.PrintErrStr("Could not add API key to cache: " + err.Error())
					requestRespondCode(w, http.StatusInternalServerError)
					return
				}
			}
		}
		requestRespond(w, newKeys)
	}
}

func APIGet(cache CacheClient[*CacheLayer]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _, err := checkUserRequest[any](r, cache)
		if err != nil {
			Printing.PrintErrStr("Could not get API: " + err.Error())
			requestRespondCode(w, http.StatusForbidden)
			return
		}
		keysInfo, err := cache.getAPIKeyInfo()
		if err != nil {
			Printing.PrintErrStr("Could not get API keys: " + err.Error())
			requestRespondCode(w, http.StatusInternalServerError)
			return
		}
		requestRespond(w, keysInfo)
	}
}
