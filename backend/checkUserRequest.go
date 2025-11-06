package main

import (
	"errors"
	"net/http"
	"reflect"
)

func formatUserRequest[ReturnType any](r *http.Request, jwt JWTService) (*ReturnType, error) {
	err := jwt.ReadAndValidateJWT(r)
	if err != nil {
		return nil, errors.New("Could not parse JWT: " + err.Error())
	}
	// Check if ReturnType is any/interface{}
	var zero ReturnType
	if reflect.TypeOf((*ReturnType)(nil)).Elem() == reflect.TypeOf((*any)(nil)).Elem() {
		return &zero, nil
	}

	requestGroup, err := requestReceived[ReturnType](r)
	return requestGroup, err
}
