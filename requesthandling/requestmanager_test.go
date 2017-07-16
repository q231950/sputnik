package requesthandling

import (
	"fmt"
	"testing"
	"time"

	"github.com/q231950/sputnik/keymanager"
	keymanagerMock "github.com/q231950/sputnik/keymanager/mocks"
)

func ExampleRequestManagerT() {
	keyManager := keymanager.New()
	containerID := "iCloud.com.elbedev.bishcommunity"
	config := RequestConfig{Version: "1", ContainerID: containerID}
	database := "public"
	requestManager := New(config, &keyManager, database)
	fmt.Println(requestManager)
	// Output: something
}

func TestPostRequest(t *testing.T) {
	keyManager := keymanagerMock.MockKeyManager{}
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
	keyManager := keymanagerMock.MockKeyManager{}
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
	keyManager := keymanagerMock.MockKeyManager{}
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
	keyManager := keymanagerMock.MockKeyManager{}
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
	keyManager := keymanagerMock.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	database := "public"
	requestManager := New(config, &keyManager, database)
	message := requestManager.message("date", "body", "service url")
	if message != "date:body:service url" {
		t.Errorf("The request payload needs to be properly formatted")
	}
}
