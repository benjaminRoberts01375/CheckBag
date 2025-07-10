package models

import (
	"os"
	"strconv"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

type config struct {
	DevMode bool `json:"dev_mode"`
}

var Config config

func Setup() {
	Coms.ReadExternalConfig("config.json", &Config)

	devMode, err := strconv.ParseBool(os.Getenv("DEV_MODE"))
	if err != nil {
		panic("Failed to parse DEV_MODE: " + err.Error())
	}
	Config.DevMode = devMode
}
