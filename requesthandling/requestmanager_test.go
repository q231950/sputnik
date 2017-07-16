package requesthandling

import (
	"fmt"
	"testing"
	"time"

	mocks "github.com/q231950/sputnik/keymanager/mocks"
)

// This Example shows how to create a request manager.
//
// A request manager requires a keymanager for handling authentication as well as a valid configuration. It currently also needs to know which database to talk to.
func ExampleRequestManager() {
	keyManager := mocks.MockKeyManager{}
	containerID := "iCloud.com.mycontainer"
	config := RequestConfig{Version: "1", ContainerID: containerID}
	database := "public"
	requestManager := New(config, &keyManager, database)
	fmt.Printf("Container:%s, Version:%s",
		requestManager.Config.ContainerID,
		requestManager.Config.Version)
	// Output: Container:iCloud.com.mycontainer, Version:1
}

func TestPostRequest(t *testing.T) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	database := "public"
	requestManager := New(config, &keyManager, database)
	request, err := requestManager.PostRequest("records/modify", ``)

	if request == nil {
		t.Errorf("The Post Request must not be nil")
	}

	if err != nil {
		t.Errorf("A Post Request must not result in error")
	}
}

func TestNewRequestManager(t *testing.T) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	database := "database"
	requestManager := New(config, keyManager, database)
	if requestManager.keyManager != keyManager {
		t.Errorf("A Request Manager's key manager should be the same that was used at initialisation")
	}

	if requestManager.database != "database" {
		t.Errorf("A Request Manager's database should not change after initialisation")
	}
}

func TestPostRequest2(t *testing.T) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	database := "public"
	requestManager := New(config, &keyManager, database)
	request, err := requestManager.PostRequest("records/modify", ``)

	if request == nil {
		t.Errorf("The Post Request must not be nil")
	}

	if err != nil {
		t.Errorf("A Post Request must not result in error")
	}
}

func TestRequest(t *testing.T) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	database := "public"
	requestManager := New(config, &keyManager, database)
	request, _ := requestManager.Request("/some/subpath", "GET", `{"key":"value", "keys":["value1", "value2"]}`)
	dateString := request.Header.Get("X-Apple-CloudKit-Request-ISO8601Date")
	if request == nil {
		t.Errorf("The Request must not be nil")
	}

	expectedTime := time.Now().UTC()
	roundedExpectedTime := expectedTime.Round(time.Minute)

	actualTime, _ := time.Parse(time.RFC3339, dateString)
	roundedTime := actualTime.Round(time.Minute)

	if !roundedExpectedTime.Equal(roundedTime) {
		t.Errorf("The date parameter must not differ by more than a minute from now")
	}
}

func TestPayloadFormat(t *testing.T) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	database := "public"
	requestManager := New(config, &keyManager, database)
	message := requestManager.message("date", "body", "service url")
	if message != "date:body:service url" {
		t.Errorf("The request payload needs to be properly formatted")
	}
}
