package models

import (
	"os"
	"strconv"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

type config struct {
	DevMode         bool `json:"dev_mode"`
	ServiceIDLength int  `json:"service_id_length"`
}

var Config config

func Setup() {
	Coms.ReadExternalConfig("config.json", &Config)

	devMode, err := strconv.ParseBool(os.Getenv("DEV_MODE"))
	if err != nil {
		panic("Failed to parse DEV_MODE: " + err.Error())
	}
	serviceIDLen, err := strconv.Atoi(os.Getenv("SERVICE_ID_LENGTH"))
	if err != nil {
		panic("Failed to parse SERVICE_ID_LENGTH: " + err.Error())
	}
	Config.ServiceIDLength = serviceIDLen
	Config.DevMode = devMode
}
