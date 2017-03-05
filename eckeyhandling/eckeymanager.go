package eckeyhandling

import (
	"bytes"
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"sync"
)

type KeyManager struct {
}

func (k *KeyManager) KeyId() string {
	return "abc"
}

func (k *KeyManager) PublicKey() *ecdsa.PublicKey {
	fmt.Println("get the public key from me")

	pemString := PrivatePublicKeyWriter()
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

func PrivatePublicKeyWriter() string {
	ecKeyPath, pathErr := ecKeyPath()
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

func ECKeyExists() bool {
	ecKeyPath, err := ecKeyPath()
	if err != nil {
		fmt.Println(err)
	}

	file, openError := os.Open(ecKeyPath)
	if openError != nil {
		fmt.Println(openError)
	}

	return file != nil && openError == nil
}

func ecKeyPath() (string, error) {
	secretsFolder, err := SecretsFolder()
	return secretsFolder + "/eckey.pem", err
}

func ECKey() string {
	ecKeyPath, pathErr := ecKeyPath()
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

func CreateECKey() error {
	ecKeyPath, pathErr := ecKeyPath()
	if pathErr != nil {
		log.Fatal(pathErr)
	}

	command := exec.Command("openssl", "ecparam", "-name", "prime256v1", "-genkey", "-noout", "-out", ecKeyPath)
	err := command.Start()
	if err != nil {
		log.Fatal(err)
	}

	err = command.Wait()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
	}

	return err
}

func SecretsFolder() (string, error) {
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
