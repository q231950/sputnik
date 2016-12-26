package requesthandling

import (
	"net/http"
)

type RequestManager struct {
}

func (r *RequestManager) PingRequest() (*http.Request, error) {
	return http.NewRequest("GET", "https://elbedev.com", nil)
}
