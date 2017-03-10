package requesthandling

import (
	"github.com/q231950/sputnik/keymanager/mocks"
	"testing"
	"time"
)

func TestPingRequest(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	requestManager := CloudkitRequestManager{keyManager}
	request, err := requestManager.PingRequest()

	if request == nil {
		t.Errorf("The Ping Request must not be nil")
	}

	if err != nil {
		t.Errorf("A Ping Request must not result in error")
	}
}

func TestPingRequestDateParameterIsInPerimeter(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	requestManager := CloudkitRequestManager{keyManager}
	request, _ := requestManager.PingRequest()
	dateString := request.Header.Get("X-Apple-CloudKit-Request-ISO8601Date")

	expectedTime := time.Now()
	roundedExpectedTime := expectedTime.Round(time.Minute)

	actualTime, _ := time.Parse("2006-01-02T15:04:05MST-0700", dateString)
	roundedTime := actualTime.Round(time.Minute)

	if !roundedExpectedTime.Equal(roundedTime) {
		t.Errorf("The date parameter must not differ by more than a minute from now")
	}
}

func TestPayloadFormat(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	requestManager := CloudkitRequestManager{keyManager}
	payload := requestManager.payload("date", "body", "service url")
	if payload != "date:body:service url" {
		t.Errorf("The request payload needs to be properly formatted")
	}
}
