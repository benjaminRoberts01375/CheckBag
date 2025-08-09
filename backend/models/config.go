package models

import (
	"os"
	"strconv"
)

type config struct {
	DevMode         bool `json:"dev_mode"`
	ServiceIDLength int  `json:"service_id_length"`
	LaunchPort      int  `json:"launch_port"`
}

var ModelsConfig config

func Setup() {
	devMode, err := strconv.ParseBool(os.Getenv("DEV_MODE"))
	if err != nil {
		panic("Failed to parse DEV_MODE: " + err.Error())
	}
	serviceIDLen, err := strconv.Atoi(os.Getenv("SERVICE_ID_LENGTH"))
	if err != nil {
		panic("Failed to parse SERVICE_ID_LENGTH: " + err.Error())
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		panic("Failed to parse PORT: " + err.Error())
	}
	ModelsConfig.ServiceIDLength = serviceIDLen
	ModelsConfig.DevMode = devMode
	ModelsConfig.LaunchPort = port
}

func (config *config) FormatPort() string {
	return ":" + strconv.Itoa(int(config.LaunchPort))
}
