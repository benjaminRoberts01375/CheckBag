package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/gorilla/websocket"
)

// Global upgrader for WebSocket connections.
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Attempts act as a proxy server for external requests to internal services. Ex. dev.benlab.us/my/stuff -> 192.168.0.50:8154/my/stuff
func requestForwarding(w http.ResponseWriter, r *http.Request) {
	requestedService := serviceLinks.GetServiceFromExternalURL(r.Host)
	if requestedService == nil {
		Coms.PrintErrStr("No service found for external URL: " + r.Host)
		requestRespondCode(w, http.StatusNotFound)
		return
	}
	internalAddress := requestedService.InternalAddress
	if internalAddress[len(internalAddress)-1] == '/' { // Remove trailing slash
		internalAddress = internalAddress[:len(internalAddress)-1]
	}
	path := r.PathValue("path")
	if len(path) > 0 && path[0] != '/' { // Add leading slash
		path = "/" + path
	}
	internalAddress += path

	if websocket.IsWebSocketUpgrade(r) {
		Coms.Println("WebSocket upgrade request received.")
		websocketProxy(w, r, internalAddress)
		return
	}
	internalAddress = "http://" + internalAddress
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

// websocketProxy handles the WebSocket connection upgrade and message forwarding.
func websocketProxy(w http.ResponseWriter, r *http.Request, baseInternalURL string) {
	requestedService := serviceLinks.GetServiceFromExternalURL(r.Host)
	if requestedService == nil {
		Coms.PrintErrStr("No service found for external URL: " + r.Host)
		requestRespondCode(w, http.StatusNotFound)
		return
	}

	internalURL, err := url.Parse("ws://" + baseInternalURL)
	if err != nil {
		Coms.PrintErrStr("Invalid internal service address: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}

	Coms.Println("Attempting to dial internal WebSocket: ", internalURL)

	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Coms.PrintErrStr("Failed to upgrade client connection: " + err.Error())
		return
	}
	defer clientConn.Close()

	// Define a map of headers to exclude. Using a map for efficient lookups.
	// These are headers that websocket.DefaultDialer will set automatically or are specific to the client-side upgrade.
	excludedHeaders := map[string]struct{}{
		"Upgrade":                  {},
		"Connection":               {},
		"Sec-Websocket-Version":    {},
		"Sec-Websocket-Key":        {},
		"Sec-Websocket-Extensions": {},
		"Host":                     {}, // Exclude the original Host header as we are setting it explicitly below
	}

	// Create a new header map and copy only non-WebSocket-specific headers.
	dialHeaders := make(http.Header)

	// Explicitly set the Host header for the internal connection to match the internal service's host.
	// This is crucial for backends that are sensitive to the Host header.
	if internalURL.Host != "" {
		dialHeaders.Set("Host", internalURL.Host)
	}

	for name, values := range r.Header {
		canonicalName := http.CanonicalHeaderKey(name)
		if _, exists := excludedHeaders[canonicalName]; exists {
			continue // Skip this header
		}
		for _, value := range values {
			dialHeaders.Add(name, value)
		}
	}

	internalConn, _, err := websocket.DefaultDialer.Dial(internalURL.String(), dialHeaders)
	if err != nil {
		Coms.PrintErrStr("Failed to connect to internal WebSocket service (" + internalURL.String() + "): " + err.Error())
		clientConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal service unavailable"))
		return
	}
	defer internalConn.Close()

	Coms.Println("Successfully established WebSocket connection to internal service.")

	errChan := make(chan error, 2)

	go func() {
		defer func() {
			errChan <- nil
		}()
		for {
			messageType, p, err := clientConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					Coms.Println("Client WebSocket connection closed normally.")
				} else {
					Coms.PrintErrStr("Error reading from client WebSocket: " + err.Error())
				}
				return
			}
			if err := internalConn.WriteMessage(messageType, p); err != nil {
				Coms.PrintErrStr("Error writing to internal WebSocket: " + err.Error())
				return
			}
		}
	}()

	go func() {
		defer func() {
			errChan <- nil
		}()
		for {
			messageType, p, err := internalConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
					Coms.Println("Internal WebSocket connection closed normally.")
				} else {
					Coms.PrintErrStr("Error reading from internal WebSocket: " + err.Error())
				}
				return
			}
			if err := clientConn.WriteMessage(messageType, p); err != nil {
				Coms.PrintErrStr("Error writing to client WebSocket: " + err.Error())
				return
			}
		}
	}()

	<-errChan
	<-errChan
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
	w.WriteHeader(http.StatusInternalServerError)
	return requestRespond(w, []byte(err.Error()), headers...)
}
