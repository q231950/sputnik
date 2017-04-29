package sputnik

import (
	"fmt"
	"log"

	"github.com/q231950/sputnik/keymanager"
	"github.com/q231950/sputnik/requesthandling"
)

func Hello() {
	fmt.Println("This is спутник.")

	keyManager := keymanager.New()
	subpath := "records/modify"
	database := "public"
	requestManager := requesthandling.New(&keyManager, database, subpath)
	request, err := requestManager.PingRequest()

	if err == nil {
		fmt.Println(request)
	} else {
		log.Fatal("Failed to create ping request")
	}
}
