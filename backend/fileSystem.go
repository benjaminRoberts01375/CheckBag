package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strconv"

	Printing "github.com/benjaminRoberts01375/CheckBag/backend/logging"
)

func init() {
	fileSystem.Setup()
}

var fileSystem = FileSystem{
	BasePath: "/CheckBag",
	UserData: "userdata.txt",
	Services: "services.json",
}

type FileSystem struct {
	BasePath string `json:"base_path"`
	UserData string `json:"user_data"`
	Services string `json:"services"`
}

func (fs *FileSystem) Setup() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	fs.BasePath = filepath.Join(homeDir, fs.BasePath)
	if err := os.MkdirAll(fs.BasePath, os.ModePerm); err != nil {
		panic(err)
	}
}

type FileServiceLinks struct {
	Version int `json:"version"`
	ServiceLinks
}

const FileServiceLinksVersion = 1

func (fs *FileSystem) Write(path string, data string) error {
	newFile, err := os.Create(path)
	if err != nil {
		return errors.New("Could not create file: " + err.Error())
	}
	defer newFile.Close()
	_, err = newFile.WriteString(data)
	if err != nil {
		return errors.New("Could not write to file: " + err.Error())
	}
	return nil
}

func (fs *FileSystem) GetUserData() (string, error) {
	data, err := os.ReadFile(filepath.Join(fs.BasePath, fs.UserData))
	return string(data), err
}

func (fs *FileSystem) SetUserData(data string) error {
	return fs.Write(filepath.Join(fs.BasePath, fs.UserData), data)
}

func (fs *FileSystem) SetServices(services ServiceLinks) error {
	updatedServices := FileServiceLinks{
		FileServiceLinksVersion, // Version
		services,                // Existing data
	}

	data, err := json.Marshal(updatedServices)
	if err != nil {
		return errors.New("Could not marshal services: " + err.Error())
	}
	return fs.Write(filepath.Join(fs.BasePath, fs.Services), string(data))
}

func (fs *FileSystem) GetServices() (ServiceLinks, error) {
	// Read data - Check if the file exists
	data, err := os.ReadFile(filepath.Join(fs.BasePath, fs.Services))
	if err != nil {
		return ServiceLinks{}, errors.New("Could not read services: " + err.Error())
	}

	// Read just the version from the file
	type ServiceVersion struct {
		Version int `json:"version"`
	}
	var version ServiceVersion
	err = json.Unmarshal(data, &version)
	if err != nil { // Assume v0 services
		version = ServiceVersion{0}
	}

	// Handle version differences
	switch version.Version {
	case 0:
		Printing.Println("v0 services detected, migrating")
		var services ServiceLinksV1
		err = json.Unmarshal(data, &services)
		if err != nil {
			return ServiceLinks{}, errors.New("Could not unmarshal v0 services: " + err.Error())
		}
		migratedServices := services.Migrate()
		fs.SetServices(migratedServices) // Overwrite old services with new version
		return migratedServices, nil
	case 1:
		Printing.Println("Services are up-to-date (v1)")
		var services FileServiceLinks
		err = json.Unmarshal(data, &services)
		if err != nil {
			return ServiceLinks{}, errors.New("Could not unmarshal v1 services: " + err.Error())
		}
		return services.ServiceLinks, nil
	}

	return ServiceLinks{}, errors.New("Unknown services version: " + strconv.Itoa(version.Version))
}
