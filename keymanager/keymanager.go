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
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
)

type KeyManager struct {
	pemFileName string
	derFileName string
}

func New() KeyManager {
	return KeyManager{pemFileName: "eckey.pem", derFileName: "cert.der"}
}

func (k *KeyManager) KeyId() string {
	return "abc"
}

func (k *KeyManager) PrivateKey() *ecdsa.PrivateKey {
	fmt.Println("get the private key from me")

	path, err := k.derFilePath()
	if err != nil {
		log.Fatal("No der file found. create one by `sputnik eckey create`")
	}
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Failed to read der file at path: ", path)
	}

	privateKey, err := x509.ParseECPrivateKey(bytes)
	if err != nil {
		log.Fatal("Failed to parse private ec key from pem: ", err)
	}

	return privateKey
}

func (k *KeyManager) PublicKey() *ecdsa.PublicKey {
	fmt.Println("get the public key from me")

	pemString := k.PrivatePublicKeyWriter()
	pemData := []byte(pemString)
	block, rest := pem.Decode(pemData)
	if block == nil || block.Type != "PUBLIC KEY" {
		log.Fatal("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Got a %T, with remaining data: %q", pub, rest)

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		fmt.Println("pub is of type RSA:", pub)
	case *dsa.PublicKey:
		fmt.Println("pub is of type DSA:", pub)
	case *ecdsa.PublicKey:
		fmt.Println("pub is of type ECDSA:", pub)
		return pub
	default:
		panic("unknown type of public key")
	}

	return nil
}

func (k *KeyManager) PrivatePublicKeyWriter() string {
	ecKeyPath, pathErr := k.pemFilePath()
	if pathErr != nil {
		log.Fatal(pathErr)
	}

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

func (k *KeyManager) ECKeyExists() bool {
	ecKeyPath, err := k.pemFilePath()
	if err != nil {
		fmt.Println(err)
	}

	file, openError := os.Open(ecKeyPath)
	if openError != nil {
		fmt.Println(openError)
	}

	return file != nil && openError == nil
}

func (k *KeyManager) derFilePath() (string, error) {
	secretsFolder, err := k.SecretsFolder()
	return secretsFolder + "/" + k.derFileName, err
}

func (k *KeyManager) pemFilePath() (string, error) {
	secretsFolder, err := k.SecretsFolder()
	return secretsFolder + "/" + k.pemFileName, err
}

func (k *KeyManager) ECKey() string {
	ecKeyPath, pathErr := k.pemFilePath()
	if pathErr != nil {
		log.Fatal(pathErr)
	}

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

func (k *KeyManager) CreateSigningIdentity() error {
	pemFilePath, pathErr := k.pemFilePath()
	if pathErr != nil {
		log.Fatal(pathErr)
	}

	// create pem
	command := exec.Command("openssl", "ecparam", "-name", "prime256v1", "-genkey", "-noout", "-out", pemFilePath)
	err := command.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = command.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	// create der from pem
	// openssl ec -outform der -in eckey.pem -out cert.der
	derFilePath, pathErr := k.derFilePath()
	if pathErr != nil {
		log.Fatal(pathErr)
	}

	command = exec.Command("openssl", "ec", "der", "-in", k.pemFileName, "-out", k.derFileName, derFilePath)
	err = command.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = command.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return err
}

func (k *KeyManager) SecretsFolder() (string, error) {
	homeDir, err := homeDir()
	if err != nil {
		// can't find home directory
		fmt.Println(err)
	}

	components := []string{homeDir, ".sputnik", "secrets"}
	configFolder := strings.Join(components, "/")

	file, err := os.Open(configFolder)
	if err != nil {
		// secrets folder doesn't exist
		fmt.Println(err)
		path, createErr := createSecretsFolder(configFolder)
		return path, createErr
	}

	return file.Name(), err
}

func createSecretsFolder(in string) (string, error) {
	mode := int(0755)
	return in, os.MkdirAll(in, os.FileMode(mode))
}

func homeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		// can't get current user
		fmt.Println(err)
	}

	return usr.HomeDir, err
}
