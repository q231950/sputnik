package requesthandling

import (
	"testing"
	"unicode/utf8"
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

func TestPingRequestDateParameter(t *testing.T) {
	requestManager := RequestManager{}
	request, _ := requestManager.PingRequest()
	dateString := request.Header.Get("X-Apple-CloudKit-Request-ISO8601Date")
	if utf8.RuneCountInString(dateString) <= 0 {
		t.Errorf("The Ping Request's date header is required to be set")
	}
}
