package keymanager

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"

	log "github.com/Sirupsen/logrus"
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
	pemFileName        string
	derFileName        string
	keyIDFileName      string
	inMemoryKeyID      string
	inMemoryPrivateKey *ecdsa.PrivateKey
	inMemoryPublicKey  *ecdsa.PublicKey
}

// New returns a CloudKitKeyManager
func New() CloudKitKeyManager {
	return CloudKitKeyManager{pemFileName: "eckey.pem", derFileName: "cert.der", keyIDFileName: "keyid.txt"}
}

// KeyID looks up the CloudKit Key ID
func (c *CloudKitKeyManager) KeyID() string {
	keyID := os.Getenv("SPUTNIK_CLOUDKIT_KEYID")
	if len(keyID) <= 0 {
		// no KeyID found in environment variables
		keyIDFromFile, err := c.storedKeyID()
		if err != nil {
			// no KeyID stored, none in env var, so it's missing
			log.Warn("No Cloudkit KeyID specified. Please either provide one by `sputnik keyid store <your KeyID>`.")
			log.Fatal(err)
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
	log.Info("Key ID", c, c.inMemoryKeyID)
	if len(c.inMemoryKeyID) > 0 {
		log.Info("Returning in-memory KeyID")
		return c.inMemoryKeyID, nil
	}

	log.Info("Need to read Key ID from file")
	path := c.keyIDFilePath()
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	c.inMemoryKeyID = string(keyBytes)
	log.Warn(c, c.inMemoryKeyID)
	return c.inMemoryKeyID, nil
}

// PrivateKey returns the x509 private key that was generated when creating the signing identity
func (c *CloudKitKeyManager) PrivateKey() *ecdsa.PrivateKey {
	if c.inMemoryPrivateKey != nil {
		log.WithFields(log.Fields{}).Info("Returning in memory private key")
		return c.inMemoryPrivateKey
	}
	path := c.derFilePath()
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("No der file found. create one by `sputnik eckey create`")
	}

	privateKey, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		log.Fatal("Failed to parse private ec key from pem: ", err)
	}

	c.inMemoryPrivateKey = privateKey
	return c.inMemoryPrivateKey
}

// PublicKey returns the public key that was generated when creating the signing identity
func (c *CloudKitKeyManager) PublicKey() *ecdsa.PublicKey {
	if c.inMemoryPublicKey != nil {
		log.Warn("Returning in memory public key")
		return c.inMemoryPublicKey
	}

	pemString := c.PrivatePublicKeyWriter()
	pemData := []byte(pemString)
	block, _ := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		fmt.Println("pub is of type RSA:", pub)
	case *dsa.PublicKey:
		fmt.Println("pub is of type DSA:", pub)
	case *ecdsa.PublicKey:
		c.inMemoryPublicKey = pub
		return pub
	default:
		panic("unknown type of public key")
	}

	return nil
}

// PrivatePublicKeyWriter should be named differently. It reads and returns the PEM encoded signing identity
func (c *CloudKitKeyManager) PrivatePublicKeyWriter() string {
	ecKeyPath := c.pemFilePath()

	command := exec.Command("openssl", "ec", "-in", ecKeyPath, "-pubout")
	bytes, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

// SigningIdentityExists checks if a signing identity has been created
func (c *CloudKitKeyManager) SigningIdentityExists() bool {
	ecKeyPath := c.pemFilePath()

	file, openError := os.Open(ecKeyPath)
	defer file.Close()

	if openError != nil {
		fmt.Println(openError)
	}

	exists := file != nil && openError == nil
	return exists
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
		log.Fatal(err)
	}
	return string(bytes)
}

// CreateSigningIdentity creates a new signing identity
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
	removePemCommand := exec.Command("rm", c.pemFilePath())
	err := removePemCommand.Run()
	if err != nil {
		log.Info("Unable to remove PEM", err)
	}

	removeDerCommand := exec.Command("rm", c.derFilePath())
	err = removeDerCommand.Run()
	if err != nil {
		log.Info("Unable to remove DER", err)
	}

	removeKeyIDCommand := exec.Command("rm", c.keyIDFilePath())
	err = removeKeyIDCommand.Run()
	if err != nil {
		log.Info("Unable to remove key ID", err)
	}

	return err
}

// createPemEncodedCertificate creates the PEM encoded certificate and stores it
func (c *CloudKitKeyManager) createPemEncodedCertificate() error {
	fmt.Println("Creating PEM...")
	pemFilePath := c.pemFilePath()

	command := exec.Command("openssl", "ecparam", "-name", "prime256v1", "-genkey", "-noout", "-out", pemFilePath)
	err := command.Start()
	if err != nil {
		log.Fatal("Failed to create pem encoded certificate", err)
	}

	err = command.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	fmt.Println("Done creating PEM")

	return err
}

// createDerEncodedCertificate converts the PEM encoded signing identity to DER and stores it
func (c *CloudKitKeyManager) createDerEncodedCertificate() error {
	fmt.Println("Creating DER...", c.pemFileName, c.derFileName, c.SecretsFolder())
	inPathPem := c.SecretsFolder() + "/" + c.pemFileName
	outPathDer := c.SecretsFolder() + "/" + c.derFileName
	command := exec.Command("openssl", "ec", "-outform", "der", "-in", inPathPem, "-out", outPathDer)

	err := command.Run()
	if err != nil {
		log.Fatal("Failed to create der encoded certificate", err)
	}

	return err
}

// SecretsFolder returns the path to Sputnik's secrets folder
func (c *CloudKitKeyManager) SecretsFolder() string {
	homeDir := homeDir()

	components := []string{homeDir, ".sputnik", "secrets"}
	configFolder := strings.Join(components, "/")

	file, err := os.Open(configFolder)
	defer file.Close()

	if err != nil {
		// secrets folder doesn't exist yet, so create it
		path, createErr := createSecretsFolder(configFolder)
		if createErr != nil {
			fmt.Println(createErr)
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
		fmt.Println(err)
	}

	return usr.HomeDir
}
