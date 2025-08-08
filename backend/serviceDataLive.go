package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	Coms "github.com/benjaminRoberts01375/Go-Communicate"
)

var serviceUpdater = ServiceUpdater{
	Subscribers: make([]chan []byte, 0),
	Remove:      make(chan chan []byte, 5),
	Add:         make(chan chan []byte, 5),
	Update:      make(chan RequestAnalyticData, 5),
}

type ServiceUpdater struct {
	Subscribers []chan []byte `json:"subscribers"`
	Remove      chan chan []byte
	Add         chan chan []byte
	Update      chan RequestAnalyticData
}

type RequestAnalyticData struct {
	ServiceID    string `json:"service"`
	Resource     string `json:"resource"`
	Country      string `json:"country"`
	IP           string `json:"ip"`
	ResponseCode int    `json:"response_code"`
}

func getServiceDataLive(w http.ResponseWriter, r *http.Request) {
	Coms.Println("Received SSE request")
	_, _, err := checkUserRequest[any](r)
	if err != nil {
		Coms.PrintErrStr("Could not verify user for analytic data: " + err.Error())
		Coms.ExternalPostRespondCode(http.StatusForbidden, w)
		return
	}
	Coms.Println("Verified user")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		Coms.PrintErrStr("Could not get flusher for SSE")
		return
	}
	Coms.Println("Got flusher")
	updateChannel := make(chan []byte, 10)
	serviceUpdater.Add <- updateChannel
	Coms.Println("Added update channel")
	for {
		select { // Wait on either a new update or a close
		case <-r.Context().Done():
			Coms.Println("Closing SSE connection")
			serviceUpdater.Remove <- updateChannel
			return
		case update := <-updateChannel:
			fmt.Fprintf(w, "data: %s\n\n", update)
			flusher.Flush()
		}
	}
}

func (serviceUpdater *ServiceUpdater) sendServiceDataLive() {
	Coms.Println("Starting SSE")
	for {
		select {
		case removeChannel := <-serviceUpdater.Remove: // Remove subscriber
			Coms.Println("Attempting to removing subscriber: ", removeChannel)
			for i, subscriber := range serviceUpdater.Subscribers {
				if subscriber == removeChannel {
					serviceUpdater.Subscribers = append(serviceUpdater.Subscribers[:i], serviceUpdater.Subscribers[i+1:]...)
					close(removeChannel)
					Coms.Println("Removed subscriber: ", removeChannel)
					break
				}
			}

		case addChannel := <-serviceUpdater.Add: // Add subscriber
			Coms.Println("Adding subscriber: ", addChannel)
			serviceUpdater.Subscribers = append(serviceUpdater.Subscribers, addChannel)

		case update := <-serviceUpdater.Update: // Update all subscribers
			updateBytes, err := json.Marshal(update)
			if err != nil {
				Coms.PrintErrStr("Could not marshal SSE update: " + err.Error())
				continue
			}
			for _, subscriber := range serviceUpdater.Subscribers {
				subscriber <- updateBytes
			}
		}
	}
}
