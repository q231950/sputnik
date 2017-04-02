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

func TestRequest(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "containerid"}
	requestManager := CloudkitRequestManager{keyManager, config}
	request, err := requestManager.Request("/some/subpath", GET, `{"key":"value", "keys":["value1", "value2"]}`)

	if request == nil {
		t.Errorf("The Request must not be nil")
	}

	if err != nil {
		t.Errorf("A Request must not result in error")
	}
}
