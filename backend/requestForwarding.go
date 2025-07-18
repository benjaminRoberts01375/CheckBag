package main

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
	"github.com/gorilla/websocket"
)

// Attempts act as a proxy server for external requests to internal services. Ex. dev.benlab.us/my/stuff -> 192.168.0.50:8154/my/stuff
func requestForwarding(w http.ResponseWriter, r *http.Request) {
	requestedService := serviceLinks.GetServiceFromExternalURL(r.Host)
	if requestedService == nil {
		Coms.PrintErrStr("No service found for external URL: " + r.Host)
		requestRespondCode(w, http.StatusNotFound)
		return
	}
	internalAddress := requestedService.InternalAddress
	path := r.PathValue("path")
	if len(path) > 0 && path[0] != '/' { // Add leading slash
		path = "/" + path
	}
	internalAddress += path

	// Check for WebSocket upgrade
	if websocket.IsWebSocketUpgrade(r) {
		websocketProxy(w, r, internalAddress)
	} else {
		restForwarding(w, r, internalAddress)
	}
}

// Handles typical HTTP requests like GET, POST, etc.
func restForwarding(w http.ResponseWriter, r *http.Request, internalAddress string) {

	internalAddress = "http://" + internalAddress

	// Preserve query parameters for HTTP requests
	if r.URL.RawQuery != "" {
		internalAddress += "?" + r.URL.RawQuery
	}

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

	if strings.Contains(r.URL.Path, "socket.io") {
		Coms.Println("Socket.IO request - URL: " + r.URL.String())
		Coms.Println("Socket.IO request - Query: " + r.URL.RawQuery)
		Coms.Println("Socket.IO request - Internal URL: " + internalAddress)
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

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // Stop following redirects >:(
		},
	}
	proxyResponse, err := client.Do(proxyRequest)
	if err != nil {
		Coms.PrintErrStr("Error sending request: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}
	// Print the response status code
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
	if proxyResponse.StatusCode >= 300 && proxyResponse.StatusCode < 400 {
		newPath := proxyResponse.Header.Get("Location")
		if strings.HasPrefix(newPath, "http://") {
			fullURL, err := url.Parse(newPath)
			if err != nil {
				Coms.PrintErrStr("Error parsing redirect URL: " + err.Error())
				requestRespondCode(w, http.StatusInternalServerError)
				return
			}
			newPath = fullURL.Path
		}

		if len(newPath) > 0 {
			if newPath[0] != '/' {
				newPath = "/" + newPath
			}
			if newPath[len(newPath)-1] != '/' {
				newPath += "/"
			}
		}
		w.Header().Set("Location", "https://"+r.Host+newPath)
	}
	w.WriteHeader(proxyResponse.StatusCode)
	w.Write(responseBytes)
}

// websocketProxy handles the WebSocket connection upgrade and message forwarding.
// w and r are the original HTTP request and response writers
// baseInternalURL is the URL of the internal service to forward the request to. Ex. `192.168.0.50:8154/my/stuff`. Note that the path is preserved, and the protocol is assumed to be HTTP.
func websocketProxy(w http.ResponseWriter, r *http.Request, baseInternalURL string) {
	// Convert HTTP URL to WebSocket URL and preserve query parameters
	wsURL := "ws://" + baseInternalURL
	if r.URL.RawQuery != "" {
		wsURL += "?" + r.URL.RawQuery
	}

	// Create upgrader for client connection
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins for proxy
		},
	}

	// Upgrade client connection to WebSocket
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		Coms.PrintErrStr("Error upgrading client connection: " + err.Error())
		return
	}
	defer clientConn.Close()

	// Forward original headers to internal service
	headers := http.Header{}
	for name, values := range r.Header {
		// Skip connection-specific headers that shouldn't be forwarded
		switch name {
		case "Connection", "Upgrade", "Sec-Websocket-Key", "Sec-Websocket-Version", "Sec-Websocket-Extensions":
			continue
		}
		for _, value := range values { // Copy all other headers
			headers.Add(name, value)
		}
	}

	// Connect to internal WebSocket service
	Coms.Println("Attempting to connect to WebSocket: " + wsURL)
	internalConn, resp, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		Coms.PrintErrStr("Error connecting to internal WebSocket service: " + err.Error())
		if resp != nil {
			Coms.PrintErrStr("HTTP Response Status: " + resp.Status)
			if resp.Body != nil {
				bodyBytes, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				if len(bodyBytes) > 0 {
					Coms.PrintErrStr("Response Body: " + string(bodyBytes))
				}
			}
		}
		clientConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal service unavailable"))
		return
	}
	defer internalConn.Close()

	Coms.Println("WebSocket proxy established between client and " + wsURL)

	// Channel to signal when either connection closes
	done := make(chan struct{})
	go forwardSocketMessage(internalConn, clientConn, done)
	go forwardSocketMessage(clientConn, internalConn, done)

	// Wait for either connection to close
	<-done
	Coms.Println("WebSocket proxy connection closed")
}

func forwardSocketMessage(incoming *websocket.Conn, outgoing *websocket.Conn, done chan struct{}) {
	defer close(done)
	for {
		messageType, message, err := incoming.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				Coms.PrintErrStr("Client WebSocket read error: " + err.Error())
			}
			return
		}

		err = outgoing.WriteMessage(messageType, message)
		if err != nil {
			Coms.PrintErrStr("Error writing to internal WebSocket: " + err.Error())
			return
		}
	}
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
