package sputnik

import (
	"fmt"
	"log"

	requests "github.com/q231950/sputnik/requesthandling"
	"github.com/q231950/sputnik/keymanager"
)

// Request constructs a signed CloudKit request
func Request(path string, method requests.HTTPMethod, payload string) {

	keyManager := keymanager.New()
	config := requests.RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	subpath := "records/modify"
	database := "public"
	requestManager := requests.New(config, &keyManager, database, subpath)
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
