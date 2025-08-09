package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func ReadExternalConfig[Config any](fileName string, configType *Config) {
	filePath := filepath.Join("/run/secrets", fileName)
	jsonFile, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(jsonFile, &configType)
	if err != nil {
		panic(err)
	}
}
