package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/benjaminRoberts01375/CheckBag/backend/logging"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

// TimeFunc allows mocking time in tests
type TimeFunc func() time.Time

type JWTService struct {
	secret        []byte
	timeFunc      TimeFunc
	cookieName    string
	loginDuration time.Duration
}

func NewJWTService(secret string, timeGenerator TimeFunc) JWTService {
	return JWTService{
		secret:        []byte(secret),
		timeFunc:      timeGenerator,
		cookieName:    "checkbag-session-token",
		loginDuration: time.Hour*24*6 + time.Hour*12, // 6.5 days
	}
}

func (s *JWTService) GenerateJWT(duration time.Duration) (string, error) {
	now := s.timeFunc()
	claims := Claims{}
	claims.ExpiresAt = jwt.NewNumericDate(now.Add(duration))
	claims.IssuedAt = jwt.NewNumericDate(now)
	claims.NotBefore = jwt.NewNumericDate(now)
	claims.Issuer = "Backend API"
	claims.Subject = "Session Token"
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *JWTService) ValidateJWT(tokenString string) (*Claims, bool) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, false
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, false
	}
	return claims, claims.ExpiresAt.After(s.timeFunc())
}

func (s *JWTService) ReadAndValidateJWT(r *http.Request) error {
	cookie, err := r.Cookie(s.cookieName)
	if err != nil {
		return err
	} else if cookie.Value == "" {
		return errors.New("JWT is empty")
	}
	_, ok := s.ValidateJWT(cookie.Value)
	if !ok {
		return errors.New("JWT is invalid")
	}
	return nil
}

func loadJWTSecret(db AdvancedDB) JWTService {
	// Check if the JWT secret was passed in as an environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret != "" {
		logging.Println("JWT provided as an environment variable")
		return NewJWTService(jwtSecret, time.Now)
	}

	// Read the JWT secret from database
	jwtSecret, err := db.GetJWTSecret(context.Background())
	if err == nil && jwtSecret != "" {
		logging.Println("JWT secret provided in database")
		return NewJWTService(jwtSecret, time.Now)
	}

	// Generate a new JWT secret
	logging.Println("Generating JWT secret and storing in database")
	newJWT := generateRandomString(15)
	err = db.SetJWTSecret(context.Background(), newJWT)
	if err != nil {
		panic("Failed to save JWT secret to database: " + err.Error())
	}
	return NewJWTService(newJWT, time.Now)
}

func (s *JWTService) setJWT(w http.ResponseWriter) error {
	// Generate a new JWT
	token, err := s.GenerateJWT(s.loginDuration)
	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     s.cookieName,
		Value:    token,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(s.loginDuration),
		Path:     "/",
	})
	return nil
}
