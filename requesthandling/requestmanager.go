package requesthandling

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/q231950/sputnik/keymanager"
)

// The HTTPMethod defines the method of a request
type HTTPMethod string

const (
	// GET reqpresents HTTP GET
	GET HTTPMethod = "GET"
	// POST reqpresents HTTP POST
	POST = "POST"
	// PUT reqpresents HTTP PUT
	PUT = "PUT"
)

// The RequestManager interface exposes methods for creating requests
type RequestManager interface {
	PostRequest() (*http.Request, error)
	Request(path string, method HTTPMethod, payload string) (*http.Request, error)
}

// CloudkitRequestManager is the concrete implementation of RequestManager
type CloudkitRequestManager struct {
	KeyManager keymanager.KeyManager
	Config     RequestConfig
}

// Request creates a signed request with the given parameters
func (cm *CloudkitRequestManager) Request(path string, method HTTPMethod, payload string) (*http.Request, error) {
	keyID := cm.KeyManager.KeyID()
	currentDate := cm.formattedTime(time.Now())

	hashedBody := cm.HashedBody(payload)
	log.WithFields(log.Fields{
		"body": string(hashedBody)}).Info("sha256")

	encodedBody := base64.StdEncoding.EncodeToString([]byte(payload))
	log.WithFields(log.Fields{"encoded body": encodedBody}).Info("base64 of sha256")

	message := cm.message(currentDate, hashedBody, path)
	log.WithFields(log.Fields{
		"date": currentDate,
		"body": hashedBody,
		"path": path}).Info("message")

	signature := cm.SignatureForMessage([]byte(message), cm.KeyManager.PrivateKey())
	encodedSignature := string(base64.StdEncoding.EncodeToString(signature))
	log.WithFields(log.Fields{"message": encodedSignature}).Info("base64 of signed sha256")

	url := "https://api.apple-cloudkit.com" + path
	log.WithFields(log.Fields{"url": url}).Info("path")

	return cm.request(string(method), url, []byte(payload), keyID, currentDate, encodedSignature)

}

// PostRequest is a sample request, only used for experimenting purposes
func (cm *CloudkitRequestManager) PostRequest() (*http.Request, error) {
	keyID := cm.KeyManager.KeyID()
	currentDate := cm.formattedTime(time.Now())
	path := cm.fullSubpath("public/records/modify")

	body := cm.body()
	hashedBody := cm.HashedBody(body)
	log.WithFields(log.Fields{
		"body": string(hashedBody)}).Info("sha256")

	encodedBody := base64.StdEncoding.EncodeToString([]byte(body))
	log.WithFields(log.Fields{"encoded body": encodedBody}).Info("base64 of sha256")

	message := cm.message(currentDate, hashedBody, path)
	log.WithFields(log.Fields{
		"date": currentDate,
		"body": hashedBody,
		"path": path}).Info("message")

	signature := cm.SignatureForMessage([]byte(message), cm.KeyManager.PrivateKey())
	encodedSignature := string(base64.StdEncoding.EncodeToString(signature))
	log.WithFields(log.Fields{"message": encodedSignature}).Info("base64 of signed sha256")

	url := "https://api.apple-cloudkit.com" + path
	log.WithFields(log.Fields{"url": url}).Info("path")

	return cm.request("POST", url, []byte(body), keyID, currentDate, encodedSignature)
}

// request creates a request with the given parameters.
//	- method POST/GET/...
//	- body is used as body for POST requests.
//	- url the request's endpoint
//	- keyID Header parameter X-Apple-CloudKit-Request-KeyID
//	- date Header parameter X-Apple-CloudKit-Request-ISO8601Date
//	- signature Header parameter X-Apple-CloudKit-Request-SignatureV1
func (cm *CloudkitRequestManager) request(method string, url string, body []byte, keyID string, date string, signature string) (request *http.Request, err error) {
	request, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	request.Header.Set("X-Apple-CloudKit-Request-KeyID", keyID)
	request.Header.Set("X-Apple-CloudKit-Request-ISO8601Date", date)
	request.Header.Set("X-Apple-CloudKit-Request-SignatureV1", signature)
	return request, err
}

// SignatureForMessage returns the signature for the given message
func (cm *CloudkitRequestManager) SignatureForMessage(message []byte, priv *ecdsa.PrivateKey) (signature []byte) {
	rand := rand.Reader

	h := sha256.New()
	h.Write([]byte(message))

	opts := crypto.SHA256
	if priv != nil {
		signature, err := priv.Sign(rand, h.Sum(nil), opts)
		if err != nil {
			log.Error("unable to sign", err)
		}

		return signature
	}

	log.Error("Can't sign without a private key")

	return nil
}

// [path]/database/[version]/[container]/[environment]/[operation-specific subpath]
// https://api.apple-cloudkit.com/database/1/[container ID]/development/public/users/lookup/email
func (cm *CloudkitRequestManager) fullSubpath(path string) string {
	version := cm.Config.Version
	containerID := cm.Config.ContainerID
	components := []string{"/database", version, containerID, "development", path}
	return strings.Join(components, "/")
}

// https://developer.apple.com/library/content/documentation/DataManagement/Conceptual/CloutKitWebServicesReference/SettingUpWebServices/SettingUpWebServices.html#//apple_ref/doc/uid/TP40015240-CH24-SW4
func (cm *CloudkitRequestManager) url(path string) string {
	url := "https://api.apple-cloudkit.com"
	subpath := cm.fullSubpath(path)
	return strings.Join([]string{url, subpath}, "/")
}

// formattedTime returns the given time in a Cloudkit compatible formatted string
func (cm *CloudkitRequestManager) formattedTime(t time.Time) string {
	date := t.UTC().Format(time.RFC3339)
	return date
}

// http://stackoverflow.com/questions/35247436/cloudkit-server-to-server-authentication
func (cm *CloudkitRequestManager) message(date string, payload string, path string) string {
	components := []string{date, payload, path}
	message := strings.Join(components, ":")
	return message
}

// HashedBody takes the given body, hashes it, using sha256 and returns the base64 encoded result
func (cm *CloudkitRequestManager) HashedBody(body string) string {
	h := sha256.New()
	h.Write([]byte(body))
	return base64.StdEncoding.EncodeToString([]byte(h.Sum(nil)))
}

func (cm *CloudkitRequestManager) body() string {
	body := `{
    "operations": [
        {
            "operationType": "create",
            "record": {
                "recordType": "Shelve",
                "fields": {
                    "title": {
                        "value": "pure panda 🐼"
                    }
                }
            }
        }
    ]
}`
	return body
}
