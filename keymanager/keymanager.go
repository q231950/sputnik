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

const KeyIdEnvironmentVariableName = string("SPUTNIK_CLOUDKIT_KEYID")

type KeyManager interface {
	PublicKey() *ecdsa.PublicKey
	PrivateKey() *ecdsa.PrivateKey
	KeyId() string
	RemoveSigningIdentity() error
	StoreKeyId(key string) error
}

type CloudkitKeyManager struct {
	pemFileName string
	derFileName string
	keyIdFileName string
}

func New() CloudkitKeyManager {
	return CloudkitKeyManager{pemFileName: "eckey.pem", derFileName: "cert.der", keyIdFileName: "keyid.txt"}
}

func (k CloudkitKeyManager) KeyId() string {
	keyId := os.Getenv("SPUTNIK_CLOUDKIT_KEYID")
	if len(keyId) <= 0 {
		// no KeyId found in environment variables
		keyId, err := k.storedKeyId()
		if err != nil {
			// no KeyId stored, none in env var, so it's missing
			log.Warn("No Cloudkit KeyId specified. Please either provide one by `sputnik keyid store <your KeyId>`.")
			log.Fatal(err)
		}
		return keyId
	}
	return keyId
}

func (c CloudkitKeyManager) StoreKeyId(key string) error {
	path := c.keyIdFilePath()
	keyBytes := []byte(key)
	return ioutil.WriteFile(path, keyBytes, 0644)
}

func (c CloudkitKeyManager) storedKeyId() (string, error) {
	path := c.keyIdFilePath()
	keyBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(keyBytes), nil
}

func (k CloudkitKeyManager) PrivateKey() *ecdsa.PrivateKey {
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

func (k CloudkitKeyManager) PublicKey() *ecdsa.PublicKey {
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

func (k *CloudkitKeyManager) PrivatePublicKeyWriter() string {
	ecKeyPath := k.pemFilePath()

	command := exec.Command("openssl", "ec", "-in", ecKeyPath, "-pubout")
	bytes, err := command.Output()
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func (k *CloudkitKeyManager) SigningIdentityExists() bool {
	ecKeyPath := k.pemFilePath()

	file, openError := os.Open(ecKeyPath)
	if openError != nil {
		fmt.Println(openError)
	}

	return file != nil && openError == nil
}

func (k *CloudkitKeyManager) keyIdFilePath() string {
	secretsFolder := k.SecretsFolder()
	return secretsFolder + "/" + k.keyIdFileName
}

func (k *CloudkitKeyManager) derFilePath() string {
	secretsFolder := k.SecretsFolder()
	return secretsFolder + "/" + k.derFileName
}

func (k *CloudkitKeyManager) pemFilePath() string {
	secretsFolder := k.SecretsFolder()
	return secretsFolder + "/" + k.pemFileName
}

func (k *CloudkitKeyManager) ECKey() string {
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

func (k *CloudkitKeyManager) CreateSigningIdentity() error {
	err := k.createPemEncodedCertificate()
	if err != nil {
		return err
	}

	err = k.createDerEncodedCertificate()

	return err
}

func (k *CloudkitKeyManager) RemoveSigningIdentity() error {
	removePemCommand := exec.Command("rm", k.pemFilePath())
	err := removePemCommand.Run()
	if err != nil {
		return err
	}

	removeDerCommand := exec.Command("rm", k.derFilePath())
	err = removeDerCommand.Run()

	return err
}

func (k *CloudkitKeyManager) createPemEncodedCertificate() error {
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

func (k *CloudkitKeyManager) createDerEncodedCertificate() error {
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

func (k *CloudkitKeyManager) SecretsFolder() string {
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
