package jwt

import (
	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

type Config struct {
	JWTSecret string `json:"jwt_secret"`
}

var config Config

func Setup() {
	Coms.ReadExternalConfig("jwt.json", &config)
}
