package requesthandling

import (
	"net/http"
	"time"
)

type RequestManager struct {
}

func (r *RequestManager) PingRequest() (*http.Request, error) {
	request, err := http.NewRequest("Get", "https://elbedev.com", nil)
	time := time.Now()
	//2006-01-02T15:04:05MST-0700
	// timeString := time.Format("Mon Jan 2 15:04:05 -0700 MST 2006")
	timeString := time.Format("2006-01-02T15:04:05MST-0700")
	request.Header.Add("X-Apple-CloudKit-Request-ISO8601Date", timeString)

	return request, err
}
func (r *RequestManager) payload(date string, body string, service string) string {
	return ""
}
