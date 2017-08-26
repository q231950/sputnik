package requesthandling

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	mocks "github.com/q231950/sputnik/keymanager/mocks"
	"github.com/stretchr/testify/assert"
)

// This Example shows how to create a request manager.
//
// A request manager requires a keymanager for handling authentication as well as a valid configuration.
func ExampleRequestManager() {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.mycontainer", Database: "public"}
	requestManager := New(config, &keyManager)
	fmt.Printf("Container:%s, Version:%s, Database:%s",
		requestManager.Config.ContainerID,
		requestManager.Config.Version,
		requestManager.Config.Database)
	// Output: Container:iCloud.com.mycontainer, Version:1, Database:public
}

func TestNewRequestManager(t *testing.T) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev", Database: "public"}
	requestManager := New(config, keyManager)
	if requestManager.keyManager != keyManager {
		t.Errorf("A Request Manager's key manager should be the same that was used at initialisation")
	}
}

func TestPostRequest(t *testing.T) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev", Database: "public"}
	requestManager := New(config, &keyManager)
	request, _ := requestManager.PostRequest("modify", "")

	assert.Equal(t, "/database/1/iCloud.com.elbedev.shelve.dev/development/public/records/modify", request.URL.Path)
}

func samplePostRequest() (*http.Request, error) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev", Database: "public"}
	requestManager := New(config, &keyManager)
	return requestManager.request("some_operation", POST, `{"key":"value", "keys":["value1", "value2"]}`)
}

func TestRequest(t *testing.T) {
	request, err := samplePostRequest()
	assert.NotNil(t, request)
	assert.Nil(t, err)
}

func TestRequestMethod(t *testing.T) {
	request, _ := samplePostRequest()
	assert.Equal(t, "POST", request.Method)
}

func TestRequestDate(t *testing.T) {
	request, _ := samplePostRequest()
	dateString := request.Header.Get("X-Apple-CloudKit-Request-ISO8601Date")

	expectedTime := time.Now().UTC()
	roundedExpectedTime := expectedTime.Round(time.Minute)
	actualTime, _ := time.Parse(time.RFC3339, dateString)
	roundedTime := actualTime.Round(time.Minute)

	assert.Equal(t, roundedExpectedTime, roundedTime, "The date parameter must not differ by more than a minute from now")
}

func TestPayloadFormat(t *testing.T) {
	keyManager := mocks.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev", Database: "public"}
	requestManager := New(config, &keyManager)
	message := requestManager.message("date", "body", "service url")
	if message != "date:body:service url" {
		t.Errorf("The request payload needs to be properly formatted")
	}
}

func TestRequestHost(t *testing.T) {
	request, _ := samplePostRequest()
	assert.Equal(t, "api.apple-cloudkit.com", request.URL.Host)
}

func TestRequestScheme(t *testing.T) {
	request, _ := samplePostRequest()
	assert.Equal(t, "https", request.URL.Scheme)
}
