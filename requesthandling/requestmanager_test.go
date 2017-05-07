package requesthandling

import (
	"testing"
	"time"
	"strings"
	"github.com/q231950/sputnik/keymanager/mocks"
)
func TestPostRequest(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "containerid"}
	requestManager := CloudkitRequestManager{keyManager, config}
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
	requestManager := New(keyManager, "database", "subpath")
	if requestManager.keyManager != keyManager {
		t.Errorf("A Request Manager's key manager should be the same that was used at initialisation")
	}

	if requestManager.database != "database" {
		t.Errorf("A Request Manager's database should not change after initialisation")
	}

	if requestManager.operationSubpath != "subpath" {
		t.Errorf("A Request Manager's operation subpath should not change after initialisation")
	}
}

func TestPingRequest(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	requestManager := CloudkitRequestManager{keyManager, "database", "subpath"}
	request, err := requestManager.PingRequest()

	if request == nil {
		t.Errorf("The Post Request must not be nil")
	}

	if err != nil {
		t.Errorf("A Post Request must not result in error")
	}
}

func TestRequest(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}

	config := RequestConfig{Version: "1", ContainerID: "containerid"}
	requestManager := CloudkitRequestManager{keyManager, config}
	request, err := requestManager.Request("/some/subpath", GET, `{"key":"value", "keys":["value1", "value2"]}`)
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
	requestManager := CloudkitRequestManager{&keyManager, "database", "subpath"}
	message := requestManager.message("date", "body", "service url")
	if message != "date:body:service url" {
		t.Errorf("The request payload needs to be properly formatted")
	}
}

func TestUrl(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	requestManager := CloudkitRequestManager{&keyManager, "public", "users/caller"}
	url := requestManager.url()
	if !strings.HasSuffix(url, "public/users/caller") {
		t.Errorf("The url of the request manager is faulty")
	}
}

/*
func testSignatureForMessage(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	r := CloudkitRequestManager{keyManager}

	message := []byte("a message to be signed")
	priv := new(mocks.MockPrivateKey)
	signature := r.SignatureForMessage(message, priv)
	if signature != "some" {
		t.Errorf("Signature is not correct")
	}
}
*/
