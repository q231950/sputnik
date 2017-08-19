package requesthandling

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	log "github.com/apex/log"
	"github.com/q231950/sputnik/keymanager"
)

// The HTTPMethod defines the method of a request
type HTTPMethod string

const (
	// GET represents HTTP GET
	GET HTTPMethod = "GET"
	// POST represents HTTP POST
	POST = "POST"
	// PUT represents HTTP PUT
	PUT = "PUT"
)

// The RequestManager interface exposes methods for creating requests
type RequestManager interface {
	PostRequest(string, string) (*http.Request, error)
	Request(path string, method HTTPMethod, payload string) (*http.Request, error)
}

// CloudkitRequestManager is the concrete implementation of RequestManager
type CloudkitRequestManager struct {
	Config     RequestConfig
	keyManager keymanager.KeyManager
}

// New creates a new RequestManager
func New(config RequestConfig, keyManager keymanager.KeyManager) CloudkitRequestManager {
	return CloudkitRequestManager{Config: config, keyManager: keyManager}
}

// PostRequest is a convenience method for creating POST requests
func (cm CloudkitRequestManager) PostRequest(operationPath string, body string) (*http.Request, error) {
	return cm.Request(operationPath, "POST", body)
}

// Request creates a signed request with the given parameters
func (cm *CloudkitRequestManager) Request(p string, method HTTPMethod, payload string) (*http.Request, error) {
	keyID := cm.keyManager.KeyID()
	currentDate := cm.formattedTime(time.Now())
	path := cm.subpath(p)
	hashedBody := cm.HashedBody(payload)
	message := cm.message(currentDate, hashedBody, path)
	signature := cm.SignatureForMessage([]byte(message))
	encodedSignature := string(base64.StdEncoding.EncodeToString(signature))
	url := "https://api.apple-cloudkit.com" + path

	log.WithFields(log.Fields{
		"key id": keyID,
		"date":   currentDate,
		"body":   hashedBody,
		"base64 encoded signature": encodedSignature,
		"path": path}).Debug("Creating request")

	return cm.requestWithHeaders(string(method), url, []byte(payload), keyID, currentDate, encodedSignature)
}

// request creates a request with the given parameters.
//	- method POST/GET/...
//	- body is used as body for POST requests.
//	- url the request's endpoint
//	- keyID Header parameter X-Apple-CloudKit-Request-KeyID
//	- date Header parameter X-Apple-CloudKit-Request-ISO8601Date
//	- signature Header parameter X-Apple-CloudKit-Request-SignatureV1
func (cm *CloudkitRequestManager) requestWithHeaders(method string, url string, body []byte, keyID string, date string, signature string) (request *http.Request, err error) {
	request, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	request.Header.Set("X-Apple-CloudKit-Request-KeyID", keyID)
	request.Header.Set("X-Apple-CloudKit-Request-ISO8601Date", date)
	request.Header.Set("X-Apple-CloudKit-Request-SignatureV1", signature)
	return request, err
}

// SignatureForMessage returns the signature for the given message
func (cm *CloudkitRequestManager) SignatureForMessage(message []byte) (signature []byte) {
	priv := cm.keyManager.PrivateKey()
	rand := rand.Reader

	h := sha256.New()
	h.Write([]byte(message))

	opts := crypto.SHA256
	if priv != nil {
		signature, err := priv.Sign(rand, h.Sum(nil), opts)
		if err != nil {
			log.WithError(err).Error("Unable to sign message")
		}

		return signature
	}

	log.Fatal("Can't sign without a private key")

	return nil
}

func (cm *CloudkitRequestManager) subpath(path string) string {
	version := cm.Config.Version
	containerID := cm.Config.ContainerID
	components := []string{"/database", version, containerID, "development", cm.Config.Database, path}
	return strings.Join(components, "/")
}

func (cm *CloudkitRequestManager) formattedTime(t time.Time) string {
	date := t.UTC().Format(time.RFC3339)
	return date
}

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
