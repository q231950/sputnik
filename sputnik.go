package sputnik

import (
	"fmt"
	"log"

	"github.com/q231950/sputnik/keymanager"
	requests "github.com/q231950/sputnik/requesthandling"
)

// Request constructs a signed CloudKit request
func Request(path string, method requests.HTTPMethod, payload string) {

	keyManager := keymanager.New()
	requestManager := requests.CloudkitRequestManager{KeyManager: &keyManager}
	request, err := requestManager.Request(path, method, payload)

	if err == nil {
		fmt.Println(request)
	} else {
		log.Fatal("Failed to create post request")
	}
}

func Hello() {
	fmt.Println("This is спутник.")
}
