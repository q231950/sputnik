package keymanager

import (
	"bytes"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"

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
	pemFileName   string
	derFileName   string
	keyIDFileName string
}

// New returns a CloudKitKeyManager
func New() CloudKitKeyManager {
	return CloudKitKeyManager{pemFileName: "eckey.pem", derFileName: "cert.der", keyIDFileName: "keyid.txt"}
}

// KeyID looks up the CloudKit Key ID
func (k CloudKitKeyManager) KeyID() string {
	keyID := os.Getenv("SPUTNIK_CLOUDKIT_KEYID")
	if len(keyID) <= 0 {
		// no KeyID found in environment variables
		keyIDFromFile, err := k.storedKeyID()
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
func (c CloudKitKeyManager) StoreKeyID(key string) error {
	path := c.keyIdFilePath()
	keyBytes := []byte(key)
	return ioutil.WriteFile(path, keyBytes, 0644)
}

// storedKeyID looks up the Key ID in a file
func (c CloudKitKeyManager) storedKeyID() (string, error) {
	path := c.keyIdFilePath()
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(keyBytes), nil
}

// PrivateKey returns the private key that was generated when creating the signing identity
func (k CloudKitKeyManager) PrivateKey() *ecdsa.PrivateKey {
	path := k.derFilePath()
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("No der file found. create one by `sputnik eckey create`")
	}

	privateKey, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		log.Fatal("Failed to parse private ec key from pem: ", err)
	}

	return privateKey
}

// PublicKey returns the public key that was generated when creating the signing identity
func (k CloudKitKeyManager) PublicKey() *ecdsa.PublicKey {
	pemString := k.PrivatePublicKeyWriter()
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
		return pub
	default:
		panic("unknown type of public key")
	}

	return nil
}

// PrivatePublicKeyWriter should be named differently. It reads and returns the PEM encoded signing identity
func (k *CloudKitKeyManager) PrivatePublicKeyWriter() string {
	ecKeyPath := k.pemFilePath()

	command := exec.Command("openssl", "ec", "-in", ecKeyPath, "-pubout")
	bytes, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

// SigningIdentityExists checks if a signing identity has been created
func (k *CloudKitKeyManager) SigningIdentityExists() bool {
	ecKeyPath := k.pemFilePath()

	file, openError := os.Open(ecKeyPath)
	if openError != nil {
		fmt.Println(openError)
	}

	return file != nil && openError == nil
}

// keyIdFilePath represents the path to the Key ID file
func (k *CloudKitKeyManager) keyIdFilePath() string {
	secretsFolder := k.SecretsFolder()
	return secretsFolder + "/" + k.keyIDFileName
}

// derFilePath represents the path to the DER encoded certificate
func (k *CloudKitKeyManager) derFilePath() string {
	secretsFolder := k.SecretsFolder()
	return secretsFolder + "/" + k.derFileName
}

// pemFilePath represents the path to the PEM encoded certificate
func (k *CloudKitKeyManager) pemFilePath() string {
	secretsFolder := k.SecretsFolder()
	return secretsFolder + "/" + k.pemFileName
}

// ECKey should be named differently. It returns the public part of the PEM encoded signing identity
func (k *CloudKitKeyManager) ECKey() string {
	ecKeyPath := k.pemFilePath()

	command := exec.Command("openssl", "ec", "-in", ecKeyPath, "-pubout")

	var output bytes.Buffer
	var waitGroup sync.WaitGroup

	stdout, _ := command.StdoutPipe()
	writer := io.MultiWriter(os.Stdout, &output)

	waitGroup.Add(1)
	go func() {
		defer waitGroup.Done()
		io.Copy(writer, stdout)
	}()

	err := command.Run()
	if err != nil {
		log.Fatal(err)
	}
	waitGroup.Wait()
	return output.String()
}

// CreateSigningIdentity creates a new signing identity
func (k *CloudKitKeyManager) CreateSigningIdentity() error {
	err := k.createPemEncodedCertificate()
	if err != nil {
		return err
	}

	err = k.createDerEncodedCertificate()

	return err
}

// RemoveSigningIdentity removes the existing signing identity
func (k *CloudKitKeyManager) RemoveSigningIdentity() error {
	removePemCommand := exec.Command("rm", k.pemFilePath())
	err := removePemCommand.Run()
	if err != nil {
		return err
	}

	removeDerCommand := exec.Command("rm", k.derFilePath())
	err = removeDerCommand.Run()

	return err
}

// createPemEncodedCertificate creates the PEM encoded certificate and stores it
func (k *CloudKitKeyManager) createPemEncodedCertificate() error {
	fmt.Println("Creating PEM...")
	pemFilePath := k.pemFilePath()

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
func (k *CloudKitKeyManager) createDerEncodedCertificate() error {
	fmt.Println("Creating DER...", k.pemFileName, k.derFileName, k.SecretsFolder())
	inPathPem := k.SecretsFolder() + "/" + k.pemFileName
	outPathDer := k.SecretsFolder() + "/" + k.derFileName
	command := exec.Command("openssl", "ec", "-outform", "der", "-in", inPathPem, "-out", outPathDer)

	err := command.Run()
	if err != nil {
		log.Fatal("Failed to create der encoded certificate", err)
	}

	return err
}

// SecretsFolder returns the path to Sputnik's secrets folder
func (k *CloudKitKeyManager) SecretsFolder() string {
	homeDir := homeDir()

	components := []string{homeDir, ".sputnik", "secrets"}
	configFolder := strings.Join(components, "/")

	file, err := os.Open(configFolder)
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
