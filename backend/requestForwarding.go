package main

import (
	"bytes"
	"context"
	"io"
	"net/http"

	Printing "github.com/benjaminRoberts01375/Web-Tech-Stack/logging"
	"github.com/gorilla/websocket"
)

// Attempts act as a proxy server for external requests to internal services
func requestForwarding(w http.ResponseWriter, r *http.Request) {
	requestedService, err := serviceLinks.GetServiceFromExternalURL(r.Host)
	if err != nil {
		Printing.PrintErrStr("No service found for external URL \"" + r.Host + "\": " + err.Error())
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

	// Read request body
	var bodyBytes []byte
	var err error
	if r.Body != nil {
		bodyBytes, err = io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			Printing.PrintErr(err)
			return
		}
	}

	Printing.Println("Creating " + r.Method + " request to: " + internalAddress)
	proxyRequest, err := http.NewRequest(r.Method, internalAddress, bytes.NewBuffer(bodyBytes))
	if err != nil {
		Printing.PrintErrStr("Error creating new request: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}
	// Add headers to proxy request
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
		Printing.PrintErrStr("Error sending request: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}
	go analytics(r, proxyResponse.StatusCode)
	// Read proxy response
	defer proxyResponse.Body.Close()
	responseBytes, err := io.ReadAll(proxyResponse.Body)
	if err != nil {
		Printing.PrintErrStr("Error reading response body: " + err.Error())
		requestRespondCode(w, http.StatusInternalServerError)
		return
	}
	// Add proxy response headers to client response
	for name, values := range proxyResponse.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}
	// Handle redirects for client response
	if proxyResponse.StatusCode >= 300 && proxyResponse.StatusCode < 400 {
		newPath := proxyResponse.Header.Get("Location")
		// Add leading and trailing slashes
		if len(newPath) > 0 && newPath[0] != '/' {
			newPath = "/" + newPath
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
		Printing.PrintErrStr("Error upgrading client connection: " + err.Error())
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
	Printing.Println("Attempting to connect to WebSocket: " + wsURL)
	internalConn, resp, err := websocket.DefaultDialer.Dial(wsURL, headers)
	if err != nil {
		Printing.PrintErrStr("Error connecting to internal WebSocket service: " + err.Error())
		if resp != nil {
			Printing.PrintErrStr("HTTP Response Status: " + resp.Status)
			if resp.Body != nil {
				bodyBytes, _ := io.ReadAll(resp.Body)
				resp.Body.Close()
				if len(bodyBytes) > 0 {
					Printing.PrintErrStr("Response Body: " + string(bodyBytes))
				}
			}
		}
		clientConn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseInternalServerErr, "Internal service unavailable"))
		return
	}
	defer internalConn.Close()

	Printing.Println("WebSocket proxy established between client and " + wsURL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go forwardSocketMessage(ctx, internalConn, clientConn, cancel)
	go forwardSocketMessage(ctx, clientConn, internalConn, cancel)

	<-ctx.Done()
	Printing.Println("WebSocket proxy connection closed")
}

func forwardSocketMessage(ctx context.Context, incoming *websocket.Conn, outgoing *websocket.Conn, cancel context.CancelFunc) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		messageType, message, err := incoming.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				Printing.PrintErrStr("Client WebSocket read error: " + err.Error())
			}
			return
		}

		err = outgoing.WriteMessage(messageType, message)
		if err != nil {
			Printing.PrintErrStr("Error writing to internal WebSocket: " + err.Error())
			return
		}
	}
}
