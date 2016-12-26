package requesthandling

import (
	"net/http"
	"strings"
	"time"
)

type RequestManager struct {
}

func (r *RequestManager) PingRequest() (*http.Request, error) {
	request, err := http.NewRequest("Get", "https://elbedev.com", nil)
	timeString := r.formattedTime(time.Now())
	request.Header.Add("X-Apple-CloudKit-Request-ISO8601Date", timeString)

	return request, err
}

func (r *RequestManager) formattedTime(time time.Time) string {
	//2006-01-02T15:04:05MST-0700
	// timeString := time.Format("Mon Jan 2 15:04:05 -0700 MST 2006")
	return time.Format("2006-01-02T15:04:05MST-0700")
}

func (r *RequestManager) payload(date string, body string, service string) string {
	components := []string{date, body, service}
	return strings.Join(components, ":")
}
