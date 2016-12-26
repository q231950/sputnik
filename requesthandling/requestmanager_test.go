package requesthandling

import (
	"testing"
	"time"
)

func TestPingRequest(t *testing.T) {
	requestManager := RequestManager{}
	request, err := requestManager.PingRequest()

	if request == nil {
		t.Errorf("The Ping Request must not be nil")
	}

	if err != nil {
		t.Errorf("A Ping Request must not result in error")
	}
}

func TestPingRequestDateParameterIsInPerimeter(t *testing.T) {
	requestManager := RequestManager{}
	request, _ := requestManager.PingRequest()
	dateString := request.Header.Get("X-Apple-CloudKit-Request-ISO8601Date")

	expectedTime := time.Now()
	roundedExpectedTime := expectedTime.Round(time.Second)

	actualTime, _ := time.Parse("2006-01-02T15:04:05MST-0700", dateString)
	roundedTime := actualTime.Round(time.Second)

	if !roundedExpectedTime.Equal(roundedTime) {
		t.Errorf("The data parameter must not differ by more than a second from now")
	}
}
