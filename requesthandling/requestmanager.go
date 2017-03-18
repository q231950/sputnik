package requesthandling

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/q231950/sputnik/keymanager"
	"net/http"
	"strings"
	"time"
	log "github.com/Sirupsen/logrus"
)

type RequestManager interface {
	PingRequest() (*http.Request, error)
}

type CloudkitRequestManager struct {
	KeyManager keymanager.KeyManager
}

func (r *CloudkitRequestManager) PingRequest() (*http.Request, error) {
	keyId := r.KeyManager.KeyId()
	currentDate := r.formattedTime(time.Now())
	path := r.subpath()

	body := r.body()
	hashedBody := r.hashedBody(body)
	log.WithFields(log.Fields{
		"body": string(hashedBody)}).Info("sha256")

	encodedBody := base64.StdEncoding.EncodeToString([]byte(body))
	log.WithFields(log.Fields{"encoded body":encodedBody}).Info("base64 of sha256")

	message := r.message(currentDate, hashedBody, path)
	log.WithFields(log.Fields{
		"date":currentDate,
		"body":hashedBody,
		"path": path}).Info("message")

	signature := r.SignatureForMessage([]byte(message))
	encodedSignature := string(base64.StdEncoding.EncodeToString(signature))
	log.WithFields(log.Fields{"message":encodedSignature}).Info("base64 of signed sha256")

	url := "https://api.apple-cloudkit.com" + path
	log.WithFields(log.Fields{"url":url}).Info("path")

	return r.request("GET", url, []byte(encodedBody), keyId, currentDate, encodedSignature)
}

func (cm *CloudkitRequestManager) request(method string, url string, body []byte, keyId string, date string, signature string) (request *http.Request, err error) {
	request, err = http.NewRequest(method, url, bytes.NewBuffer(body))
	request.Header.Set("X-Apple-CloudKit-Request-KeyID", keyId)
	request.Header.Set("X-Apple-CloudKit-Request-ISO8601Date", date)
	request.Header.Set("X-Apple-CloudKit-Request-SignatureV1", signature)
	return request, err
}

func (cm *CloudkitRequestManager) SignatureForMessage(message []byte) (signature []byte) {
	priv := cm.KeyManager.PrivateKey()
	rand := rand.Reader

	h := sha256.New()
	h.Write([]byte(message))

	opts := crypto.SHA256
	signature, err := priv.Sign(rand, h.Sum(nil), opts)
	if err != nil {
		log.Info("unable to sign", err)
	}

	return signature
}

// [path]/database/[version]/[container]/[environment]/[operation-specific subpath]
// https://api.apple-cloudkit.com/database/1/[container ID]/development/public/users/lookup/email
func (r *CloudkitRequestManager) subpath() string {
	version := "1"
	containerId := "iCloud.com.elbedev.shelve.dev"
	// subpath := "public/records/query"
	// subpath := "public/users/lookup/email"
	subpath := "public/users/caller"

	components := []string{"/database", version, containerId, "development", subpath}
	return strings.Join(components, "/")
}

// https://developer.apple.com/library/content/documentation/DataManagement/Conceptual/CloutKitWebServicesReference/SettingUpWebServices/SettingUpWebServices.html#//apple_ref/doc/uid/TP40015240-CH24-SW4
func (r *CloudkitRequestManager) url() string {
	path := "https://api.apple-cloudkit.com"
	subpath := r.subpath()
	return strings.Join([]string{path, subpath}, "/")
}

func (r *CloudkitRequestManager) formattedTime(t time.Time) string {
	date := t.UTC().Format(time.RFC3339)
	return date
}

// http://stackoverflow.com/questions/35247436/cloudkit-server-to-server-authentication
func (r *CloudkitRequestManager) message(date string, payload string, path string) string {
	components := []string{date, payload, path}
	message := strings.Join(components, ":")
	return message
}

func (r CloudkitRequestManager) hashedBody(body string) string {
	h := sha256.New()
	h.Write([]byte(body))
	return base64.StdEncoding.EncodeToString([]byte(h.Sum(nil)))
}

func (r *CloudkitRequestManager) body() string {
	// body := `{"users":[{"emailAddress":"some@one.com"}]}`
	body := ``
	// body := `{"zoneID": "_defaultZone","query": {"recordType": "Shelve"}}`
	return body
}
