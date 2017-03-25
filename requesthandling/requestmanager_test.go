package requesthandling

import (
	"testing"
	"time"

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

func TestPostRequestDateParameterIsInPerimeter(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "containerid"}
	requestManager := CloudkitRequestManager{keyManager, config}
	request, _ := requestManager.PostRequest()
	dateString := request.Header.Get("X-Apple-CloudKit-Request-ISO8601Date")

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
	config := RequestConfig{Version: "1", ContainerID: "containerid"}
	requestManager := CloudkitRequestManager{keyManager, config}
	message := requestManager.message("date", "body", "service url")
	if message != "date:body:service url" {
		t.Errorf("The request payload needs to be properly formatted")
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
