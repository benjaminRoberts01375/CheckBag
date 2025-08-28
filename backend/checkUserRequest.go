package main

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/benjaminRoberts01375/CheckBag/backend/jwt"
)

func checkUserRequest[ReturnType any](r *http.Request) (*jwt.Claims, *ReturnType, error) {
	cookie, err := r.Cookie(jwt.CookieName)

	if err != nil {
		return nil, nil, errors.New("missing JWT cookie")
	}
	claims, ok := jwt.JWTIsValid(cookie.Value)
	if !ok {
		return nil, nil, errors.New("invalid JWT")
	}

	cacheData, err := cache.getUserSignIn(cookie.Value) // Pass JWT to cache to get user ID
	if err != nil {
		return nil, nil, errors.New("failed to get user ID from JWT: " + err.Error())
	} else if cacheData != cacheDataValid {
		return nil, nil, errors.New("invalid JWT")
	}

	// Check if ReturnType is any/interface{}
	var zero ReturnType
	if reflect.TypeOf((*ReturnType)(nil)).Elem() == reflect.TypeOf((*any)(nil)).Elem() {
		return claims, &zero, nil
	}

	requestGroup, err := requestReceived[ReturnType](r)
	return claims, requestGroup, err
}
