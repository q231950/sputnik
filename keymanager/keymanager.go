package keymanager

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"

	log "github.com/apex/log"
)

// KeyIDEnvironmentVariableName is the constant used for identifying the Key ID environment variable
const KeyIDEnvironmentVariableName = string("SPUTNIK_CLOUDKIT_KEYID")

// KeyManager exposes methods for creating, reading and removing signing identity relevant keys and IDs
type KeyManager interface {
	PublicKey() *ecdsa.PublicKey
	PrivateKey() *ecdsa.PrivateKey
	KeyID() string
	RemoveSigningIdentity() error
	StoreKeyID(key string) error
}

// CloudKitKeyManager is a concrete KeyManager
type CloudKitKeyManager struct {
	secretsFolder      string
	pemFileName        string
	derFileName        string
	keyIDFileName      string
	inMemoryKeyID      string
	inMemoryPrivateKey *ecdsa.PrivateKey
	inMemoryPublicKey  *ecdsa.PublicKey
}

// New returns a CloudKitKeyManager with the default secrets folder
// By default, a CloudKitKeyManager expects the secrets in the .sputnik folder of the home directory
func New() CloudKitKeyManager {
	homeDir := homeDir()
	components := []string{homeDir, ".sputnik", "secrets"}
	secretsFolder := strings.Join(components, "/")

	keyIDFileName := "keyid.txt"
	derFileName := "cert.der"
	pemFileName := "eckey.pem"

	return NewWithSecretsFolder(secretsFolder, keyIDFileName, derFileName, pemFileName)
}

// NewWithSecretsFolder returns a CloudKitKeyManager with a specific secrets folder
// Use this to specify a different storage location from the default
func NewWithSecretsFolder(secretsFolder string, keyIDFileName string, derFileName string, pemFileName string) CloudKitKeyManager {
	return CloudKitKeyManager{
		secretsFolder: secretsFolder,
		pemFileName:   pemFileName,
		derFileName:   derFileName,
		keyIDFileName: keyIDFileName}
}

// KeyID looks up the CloudKit Key ID
func (c *CloudKitKeyManager) KeyID() string {
	keyID := os.Getenv("SPUTNIK_CLOUDKIT_KEYID")
	if len(keyID) <= 0 {
		// no KeyID found in environment variables
		keyIDFromFile, err := c.storedKeyID()
		if err != nil {

		}
		return keyIDFromFile
	}
	return keyID
}

// StoreKeyID stores the given ID to a file in Sputnik's secrets folder
func (c *CloudKitKeyManager) StoreKeyID(key string) error {
	path := c.keyIDFilePath()
	keyBytes := []byte(key)
	return ioutil.WriteFile(path, keyBytes, 0644)
}

// storedKeyID looks up the Key ID in a file
func (c *CloudKitKeyManager) storedKeyID() (string, error) {
	if len(c.inMemoryKeyID) > 0 {
		return c.inMemoryKeyID, nil
	}

	path := c.keyIDFilePath()
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	c.inMemoryKeyID = string(keyBytes)

	return c.inMemoryKeyID, nil
}

// PrivateKey returns the x509 private key that was generated when creating the signing identity
func (c *CloudKitKeyManager) PrivateKey() *ecdsa.PrivateKey {
	if c.inMemoryPrivateKey != nil {
		return c.inMemoryPrivateKey
	}

	inPathPem := c.pemFilePath()
	command := exec.Command("openssl", "ec", "-outform", "der", "-in", inPathPem)
	bytes, _ := command.Output()

	privateKey, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		log.Error("Failed to parse private ec key from pem:")
		log.Errorf("%s", err)
		return nil
	}

	c.inMemoryPrivateKey = privateKey
	return c.inMemoryPrivateKey
}

// PublicKey returns the public key that was generated when creating the signing identity
func (c *CloudKitKeyManager) PublicKey() *ecdsa.PublicKey {
	if c.inMemoryPublicKey != nil {
		return c.inMemoryPublicKey
	}

	var err error
	var pub interface{}
	pemString := c.PublicKeyString()
	pemData := []byte(pemString)
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		err = errors.New("failed to decode PEM block containing public key")
	} else {
		pub, err = x509.ParsePKIXPublicKey(block.Bytes)

		switch pub := pub.(type) {
		case *ecdsa.PublicKey:
			c.inMemoryPublicKey = pub
			return pub
		}
	}

	if err != nil {
		log.Error("unable to parse public key from certificate:")
		log.Errorf("%s", err)
	}

	return nil
}

// PublicKeyString should be named differently. It reads and returns the public part of the PEM encoded signing identity
func (c *CloudKitKeyManager) PublicKeyString() string {
	ecKeyPath := c.pemFilePath()

	command := exec.Command("openssl", "ec", "-in", ecKeyPath, "-pubout")
	bytes, err := command.Output()
	if err != nil {
		log.Error("PublicKeyString")
		log.Error("Failed to read the public key from PEM")
		log.Errorf("%s", err)
	}
	return string(bytes)
}

// SigningIdentityExists checks if a signing identity has been created
func (c *CloudKitKeyManager) SigningIdentityExists() (bool, error) {
	ecKeyPath := c.pemFilePath()

	file, openError := os.Open(ecKeyPath)
	defer file.Close()

	return file != nil && openError == nil, openError
}

// keyIdFilePath represents the path to the Key ID file
func (c *CloudKitKeyManager) keyIDFilePath() string {
	secretsFolder := c.SecretsFolder()
	return secretsFolder + "/" + c.keyIDFileName
}

// derFilePath represents the path to the DER encoded certificate
func (c *CloudKitKeyManager) derFilePath() string {
	secretsFolder := c.SecretsFolder()
	return secretsFolder + "/" + c.derFileName
}

// pemFilePath represents the path to the PEM encoded certificate
func (c *CloudKitKeyManager) pemFilePath() string {
	secretsFolder := c.SecretsFolder()
	return secretsFolder + "/" + c.pemFileName
}

// ECKey should be named differently. It returns the public part of the PEM encoded signing identity
func (c *CloudKitKeyManager) ECKey() string {
	ecKeyPath := c.pemFilePath()

	command := exec.Command("openssl", "ec", "-in", ecKeyPath, "-pubout")

	bytes, err := command.Output()
	if err != nil {
		log.Error("ECKey")
		log.Error("Failed to read the public key from PEM")
		log.Errorf("%s", err)
	}
	return string(bytes)
}

// CreateSigningIdentity creates a new signing identity.
//
// You can paste the signing identity to your iCloud Dashboard when creating a new API Access Key.
func (c *CloudKitKeyManager) CreateSigningIdentity() error {
	err := c.createPemEncodedCertificate()
	if err != nil {
		return err
	}

	err = c.createDerEncodedCertificate()

	return err
}

// RemoveSigningIdentity removes the existing signing identity
func (c *CloudKitKeyManager) RemoveSigningIdentity() error {
	c.inMemoryPrivateKey = nil
	c.inMemoryPublicKey = nil

	removePemCommand := exec.Command("rm", c.pemFilePath())
	err := removePemCommand.Run()
	if err != nil {
		log.Error("Unable to remove PEM:")
		log.Errorf("%s", err)
	}

	removeDerCommand := exec.Command("rm", c.derFilePath())
	err = removeDerCommand.Run()
	if err != nil {
		log.Error("Unable to remove DER file")
		log.Errorf("%s", err)
	}

	removeKeyIDCommand := exec.Command("rm", c.keyIDFilePath())
	err = removeKeyIDCommand.Run()

	return err
}

// createPemEncodedCertificate creates the PEM encoded certificate and stores it
func (c *CloudKitKeyManager) createPemEncodedCertificate() error {
	log.Debug("Creating PEM...")
	pemFilePath := c.pemFilePath()

	command := exec.Command("openssl", "ecparam", "-name", "prime256v1", "-genkey", "-noout", "-out", pemFilePath)
	err := command.Start()
	if err != nil {
		log.Error("Failed to create pem encoded certificate")
		log.Errorf("%s", err)
	}

	err = command.Wait()
	if err != nil {
		log.Error("Error executing `openssl`")
		log.Errorf("%s", err)
	}

	log.Info("Done creating PEM")

	return err
}

// createDerEncodedCertificate converts the PEM encoded signing identity to DER and stores it
func (c *CloudKitKeyManager) createDerEncodedCertificate() error {
	log.Debugf("Creating DER...", c.pemFileName, c.derFileName, c.SecretsFolder())
	inPathPem := c.SecretsFolder() + "/" + c.pemFileName
	outPathDer := c.SecretsFolder() + "/" + c.derFileName
	command := exec.Command("openssl", "ec", "-outform", "der", "-in", inPathPem, "-out", outPathDer)

	err := command.Run()
	if err != nil {
		log.Error("Failed to create der encoded certificate")
		log.Errorf("%s", err)
	}

	return err
}

// SecretsFolder returns the path to Sputnik's secrets folder
func (c *CloudKitKeyManager) SecretsFolder() string {
	file, err := os.Open(c.secretsFolder)
	defer file.Close()

	if err != nil {
		// secrets folder doesn't exist yet, so create it
		path, createErr := createSecretsFolder(c.secretsFolder)
		if createErr != nil {
			log.Debugf("%s", createErr)
		}
		return path
	}

	return file.Name()
}

func createSecretsFolder(in string) (string, error) {
	mode := int(0755)
	return in, os.MkdirAll(in, os.FileMode(mode))
}

func homeDir() string {
	usr, err := user.Current()
	if err != nil {
		// can't get current user
		log.Errorf("%s", err)
	}

	return usr.HomeDir
}
