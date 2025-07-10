package main

import (
	"errors"
	"net/http"
	"reflect"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/benjaminRoberts01375/Web-Tech-Stack/jwt"
)

func checkUserRequest[ReturnType any](r *http.Request) (*jwt.Claims, string, *ReturnType, error) {
	cookie, err := r.Cookie(jwt.CookieName)

	if err != nil {
		return nil, "", nil, errors.New("missing JWT cookie")
	}
	claims, ok := jwt.JWTIsValid(cookie.Value)
	if !ok {
		return nil, "", nil, errors.New("invalid JWT")
	}

	userID, err := cache.getUserSignIn(cookie.Value) // Pass JWT to cache to get user ID
	if err != nil {
		Coms.PrintErrStr("failed to get user ID from JWT: " + err.Error())
	}

	// Check if ReturnType is any/interface{}
	var zero ReturnType
	if reflect.TypeOf((*ReturnType)(nil)).Elem() == reflect.TypeOf((*any)(nil)).Elem() {
		return claims, userID, &zero, nil
	}

	requestGroup, err := Coms.ExternalPostReceived[ReturnType](r)
	return claims, userID, requestGroup, err
}
