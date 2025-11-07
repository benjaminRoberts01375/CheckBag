package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

// TODO To be removed in CheckBag v5
func migrateFSToDB(db AdvancedDB) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}
	CheckBagPath := filepath.Join(homeDir, "CheckBag")

	ctx := context.Background()
	// Migrate user data
	userHashFile := filepath.Join(CheckBagPath, "userdata.txt")
	userHash, err := os.ReadFile(userHashFile)
	if err == nil { // User hash migration needed
		db.SetUserPasswordHash(ctx, string(userHash)) // Can panic
		os.Remove(userHashFile)
	}
	serviceDataFile := filepath.Join(CheckBagPath, "services.json")
	serviceData, err := os.ReadFile(serviceDataFile)
	if err == nil { // ServiceLink migration needed
		type FileServiceLinks struct {
			Version int `json:"version"`
			ServiceLinks
		}
		var services FileServiceLinks
		err = json.Unmarshal(serviceData, &services)
		if err != nil {
			panic("Unable to read services for migration")
		}
		err := db.setServiceLinks(ctx, services.ServiceLinks) // Can't panic
		if err != nil {
			panic("Unable to add service links to DB during migration: " + err.Error())
		}
		os.Remove(serviceDataFile)
	}
	os.RemoveAll(CheckBagPath) // You'll be missed :')
}
