package requesthandling

import (
	"testing"
	"time"

	keymanager "github.com/q231950/sputnik/keymanager/mocks"
)

type TestableRequestManager interface {
	RequestManager
	HashedBody(body string) string
}

func TestPostRequestDateParameterIsInPerimeter(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "containerid"}
	subpath := "users/caller"
	database := "public"
	requestManager := New(config, &keyManager, database, subpath)
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

func TestMessageFormat(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := RequestConfig{Version: "1", ContainerID: "containerid"}
	subpath := "users/caller"
	database := "public"
	requestManager := New(config, &keyManager, database, subpath)
	message := requestManager.message("date", "body", "service url")
	if message != "date:body:service url" {
		t.Errorf("The request payload needs to be properly formatted")
	}
}

func TestEmptyHashedBody(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := NewRequestConfig("version", "containerID")
	subpath := "users/caller"
	database := "public"
	requestManager := TestableRequestManager(&CloudkitRequestManager{config, &keyManager, database, subpath})
	body := ""
	hash := requestManager.HashedBody(body)

	if string(hash) != "47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=" {
		t.Errorf("Signature is not correct")
	}
}

func TestSignMessage(t *testing.T) {
	keyManager := keymanager.MockKeyManager{}
	config := NewRequestConfig("version", "containerID")
	subpath := "subpath"
	database := "database"
	r := New(config, &keyManager, database, subpath)
	signature := r.SignatureForMessage([]byte("message"))

	if signature == nil {
		t.Errorf("A message should be signed when a private key is available", signature)
	}
}

// func TestSignature(t *testing.T) {
//
// 	keyManager := keymanager.MockKeyManager{}
// 	config := NewRequestConfig("version", "containerID")
// 	r := CloudkitRequestManager{keyManager, config}
//   message := []byte("")
//
// 	bytes, err := ioutil.ReadFile("fixtures/test_identity.der")
// 	priv, err := x509.ParseECPrivateKey(bytes)
// 	if err != nil {
// 		t.Errorf("Incorrect test setup %s", err)
// 		t.Fail()
// 	}
// 	signature := r.SignatureForMessage(message, priv)
// 	if string(signature) != "47DEQpj8HBSa+/TImW+5JCeuQeRkm5NMpJWZG3hSuFU=" {
// 		t.Errorf("Signature is not correct")
// 	}
// }
