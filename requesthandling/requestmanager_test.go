package requesthandling

import (
	"testing"

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
