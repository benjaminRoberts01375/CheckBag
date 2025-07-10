package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}
