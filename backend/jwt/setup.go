package jwt

import (
	Config "github.com/benjaminRoberts01375/Web-Tech-Stack/config"
)

type config struct {
	JWTSecret string `json:"jwt_secret"`
}

var JWTConfig config

func Setup() {
	Config.ReadExternalConfig("jwt.json", &JWTConfig)
}
