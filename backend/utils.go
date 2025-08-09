package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

func generateRandomString(length int) string {
	// Charset is URL safe and easy to read
	const charset = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"

	stringBase := make([]byte, length)
	for i := range stringBase {
		stringBase[i] = charset[rand.Intn(len(charset))]
	}
	return string(stringBase)
}

func requestRespond(w http.ResponseWriter, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, error := w.Write(jsonData)
	return error
}

func requestRespondCode(w http.ResponseWriter, code int) error {
	w.WriteHeader(code)
	return requestRespond(w, nil)
}

func requestReceived[ReturnType any](r *http.Request) (*ReturnType, error) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var request ReturnType
	err = json.Unmarshal(body, &request)
	return &request, err
}
