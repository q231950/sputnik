package requesthandling

import (
  "testing"
)

func TestRequest(t *testing.T) {
  requestManager := RequestManager{}
  request, err := requestManager.PingRequest()

  if request == nil {
    t.Errorf("The Ping Request must not be nil")
  }

  if err != nil {
    t.Errorf("A Ping Request must not result in error")
  }
}
