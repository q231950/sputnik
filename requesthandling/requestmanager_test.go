package requesthandling

import (
	"testing"
	"time"

	"github.com/q231950/sputnik/keymanager/mocks"
)

func TestPostRequest(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	subpath := "records/modify"
	database := "public"
	requestManager := New(config, &keyManager, database, subpath)
	request, err := requestManager.PostRequest()

	if request == nil {
		t.Errorf("The Post Request must not be nil")
	}

	if err != nil {
		t.Errorf("A Post Request must not result in error")
	}
}

func TestNewRequestManager(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	subpath := "sub/path"
	database := "database"
	requestManager := New(config, keyManager, database, subpath)
	if requestManager.keyManager != keyManager {
		t.Errorf("A Request Manager's key manager should be the same that was used at initialisation")
	}

	if requestManager.database != "database" {
		t.Errorf("A Request Manager's database should not change after initialisation")
	}

	if requestManager.operationSubpath != "sub/path" {
		t.Errorf("A Request Manager's operation subpath should not change after initialisation")
	}
}

func TestPostRequest2(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	subpath := "records/modify"
	database := "public"
	requestManager := New(config, &keyManager, database, subpath)
	request, err := requestManager.PostRequest()

	if request == nil {
		t.Errorf("The Post Request must not be nil")
	}

	if err != nil {
		t.Errorf("A Post Request must not result in error")
	}
}

func TestRequest(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	subpath := "records/modify"
	database := "public"
	requestManager := New(config, &keyManager, database, subpath)
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
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "iCloud.com.elbedev.shelve.dev"}
	subpath := "records/modify"
	database := "public"
	requestManager := New(config, &keyManager, database, subpath)
	message := requestManager.message("date", "body", "service url")
	if message != "date:body:service url" {
		t.Errorf("The request payload needs to be properly formatted")
	}
}
