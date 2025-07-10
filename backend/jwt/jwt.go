package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const CookieName = "session-token"
const LoginDuration = time.Hour*24*6 + time.Hour*12 // 6 + 0.5 days

func GenerateJWT(username string, duration time.Duration) (string, error) {
	claims := Claims{
		Username: username,
	}
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.NotBefore = jwt.NewNumericDate(time.Now())
	claims.Issuer = "Backend API"
	claims.Subject = "Session Token"
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(config.JWTSecret))
}

func JWTIsValid(tokenString string) (*Claims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(config.JWTSecret), nil
	})
	if err != nil {
		return nil, false
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, false
	}
	return claims, claims.ExpiresAt.After(time.Now())
}
