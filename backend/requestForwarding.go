package main

import (
	"bytes"
	"fmt"
	"io"
	"maps"
	"net/http"
)

func requestForwarding(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\n----------")
	fmt.Println("Headers:")
	for k, v := range r.Header {
		fmt.Println(k, v)
	}
	fmt.Println("\nHost: " + r.Host)
	fmt.Println("Request Type: " + r.Method)
	fmt.Println("Path: " + r.URL.Path)
	fmt.Println("Endpoint: " + r.PathValue("path"))
	fmt.Println("Query: " + r.URL.RawQuery)
	var bodyBytes []byte
	var err error
	if r.Body != nil {
		fmt.Println("Body:")
		bodyBytes, err = io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(bodyBytes))
	}
}
