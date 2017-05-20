package sputnik

import (
	"fmt"
	"net/http"

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

// Request constructs a signed CloudKit request
func Request(path string, method HTTPMethod, payload string) (*http.Response, error) {

	keyManager := keymanager.New()
	config := requests.RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	database := "public"
	requestManager := requests.New(config, &keyManager, database, path)

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
	fmt.Println("This is спутник.")
}
