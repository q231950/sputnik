package main

import (
	"net/http"

	log "github.com/apex/log"

	"github.com/q231950/sputnik/keymanager"
	requests "github.com/q231950/sputnik/requesthandling"
)

// The HTTPMethod defines the method of a request
type HTTPMethod string

const (
	// GET represents HTTP GET
	GET HTTPMethod = "GET"
	// POST represents HTTP POST
	POST = "POST"
	// PUT represents HTTP PUT
	PUT = "PUT"
)

// Post the payload to the path
func Post(path string, payload string, containerID string) (*http.Response, error) {
	return request(path, POST, payload, containerID)
}

// Request constructs a signed CloudKit request
func request(path string, method HTTPMethod, payload string, containerID string) (*http.Response, error) {

	keyManager := keymanager.New()
	config := requests.RequestConfig{Version: "1", ContainerID: containerID, Database: "public"}
	requestManager := requests.New(config, &keyManager)

	request, err := requestManager.Request(path, string(method), payload)
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	return resp, err
}

// Hello makes sure that the earth is still spinning around the sun
func Hello() {
	log.Debug("This is спутник.")
}
