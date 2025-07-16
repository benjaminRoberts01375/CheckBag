package main

import (
	"bytes"
	"io"
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

// Attempts act as a proxy server for external requests to internal services. Ex. dev.benlab.us/my/stuff -> 192.168.0.50:8154/my/stuff
func requestForwarding(w http.ResponseWriter, r *http.Request) {
	var bodyBytes []byte
	var err error
	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			Coms.PrintErr(err)
			return
		}
	}
	requestedService := serviceLinks.GetServiceFromExternalURL(r.Host)
	if requestedService == nil {
		Coms.PrintErrStr("No service found for external URL: " + r.Host)
		requestRespondCode(w, http.StatusNotFound)
		return
	}
	internalAddress := requestedService.InternalAddress + "/" + r.PathValue("path") // Ex. 192.168.0.50:8154
	if internalAddress[len(internalAddress)-1] == '/' {
		internalAddress = internalAddress[:len(internalAddress)-1]
	}

	Coms.Println("Sending " + r.Method + " request to: " + internalAddress)
	proxyRequest, err := http.NewRequest(r.Method, internalAddress, bytes.NewBuffer(bodyBytes))
	if err != nil {
		Coms.PrintErrStr("Error creating new request: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}

	// Set headers for proxy
	for name, values := range r.Header {
		for _, value := range values {
			proxyRequest.Header.Add(name, value)
		}
	}

	client := &http.Client{}
	proxyResponse, err := client.Do(proxyRequest)
	if err != nil {
		Coms.PrintErrStr("Error sending request: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}

	// Respond to client

	defer proxyResponse.Body.Close()
	responseBytes, err := io.ReadAll(proxyResponse.Body)
	if err != nil {
		Coms.PrintErrStr("Error reading response body: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}
	for name, values := range proxyResponse.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	w.WriteHeader(proxyResponse.StatusCode)
	w.Write(responseBytes)
}

func requestRespond(w http.ResponseWriter, data []byte, headers ...http.Header) error {
	_, error := w.Write(data)
	return error
}

func requestRespondCode(w http.ResponseWriter, code int, headers ...http.Header) error {
	w.WriteHeader(code)
	return requestRespond(w, nil, headers...)
}

func requestRespondError(w http.ResponseWriter, err error, headers ...http.Header) error {
	w.WriteHeader(http.StatusInternalServerError) // Runs only if the error wasn't set already
	return requestRespond(w, []byte(err.Error()), headers...)
}
