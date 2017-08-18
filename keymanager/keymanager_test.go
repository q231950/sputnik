package keymanager

import (
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	manager := New()
	assert.NotNil(t, manager, "A newly created key manager should not be nil")
}

func TestStoredKeyID(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")
	assert.Equal(t, "abc\n", manager.KeyID(), "A stored key id should get found")
}

func TestPrivateKey(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")

	expectedD := new(big.Int)
	expectedD.SetString("57359333433306843951573675484597381433848383364258304847182053853963006392866", 10)
	assert.Equal(t, expectedD, manager.PrivateKey().D, "The private key does not match the certificate")
}

func TestPrivateKeyFromMemory(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")

	privateKey := manager.PrivateKey()
	privateKey = manager.PrivateKey()
	expectedD := new(big.Int)
	expectedD.SetString("57359333433306843951573675484597381433848383364258304847182053853963006392866", 10)
	assert.Equal(t, expectedD, privateKey.D, "The private key does not match the certificate")
}

func TestPublicKey(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")

	publicKey := manager.PublicKey()
	expectedX := new(big.Int)
	expectedX.SetString("83764337057786748884235593670888306280068598385338703097790983192099570881279", 10)
	assert.Equal(t, expectedX, publicKey.X, "The public key's X does not match the certificate")

	expectedY := new(big.Int)
	expectedY.SetString("52205868503601725522780045291819081986668283346034058232083577220213552372588", 10)
	assert.Equal(t, expectedY, publicKey.Y, "The public key's Y does not match the certificate")
}

func TestPublicKeyFromMemory(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")

	publicKey := manager.PublicKey()
	publicKey = manager.PublicKey()
	expectedX := new(big.Int)
	expectedX.SetString("83764337057786748884235593670888306280068598385338703097790983192099570881279", 10)
	assert.Equal(t, expectedX, publicKey.X, "The public key's X does not match the certificate")
}

func TestStoreKeyID(t *testing.T) {
	pathToFixtures := "./testFiles"
	manager := NewWithSecretsFolder(pathToFixtures, "keyidCreateTest.txt", "eckeyTest.pem")

	manager.StoreKeyID("key")
	assert.Equal(t, "key", manager.KeyID(), "The key id should be stored correctly")

	_ = os.RemoveAll("./testFiles")
}

func TestSigningIdentityExists(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")
	exists, _ := manager.SigningIdentityExists()
	assert.True(t, exists, "The signing identity should exist in this fixture")
}

func TestECKey(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")
	expected := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEuTDvRjrfQW7CHI68iJBpBpRkhT0J\nKMKA7vaSPvkckv9za3l1ji2L0VsFSOIvSlgpUyC96pRxcIBR/E2gqmLbbA==\n-----END PUBLIC KEY-----\n"
	assert.Equal(t, expected, manager.ECKey(), "The public key is not correct")
}

func TestCreateIdentity(t *testing.T) {
	pathToFixtures := "./testFiles"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")
	err := manager.CreateSigningIdentity()
	assert.Nil(t, err, "It should be possible to create a signing identity")

	_ = os.RemoveAll("./testFiles")
}

func TestCreateIdentityPEMError(t *testing.T) {
	pathToFixtures := "./testFiles"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "")
	err := manager.CreateSigningIdentity()
	assert.NotNil(t, err, "An error should occur when creating the signing identity fails due to a false pem file name")

	_ = os.RemoveAll("./testFiles")
}

func TestRemovesECKeyKeyWhenRemovingSigningIdentity(t *testing.T) {
	pathToFixtures := "./testFiles"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")
	_ = manager.CreateSigningIdentity()

	manager.RemoveSigningIdentity()

	assert.Equal(t, "", manager.ECKey(), "The public key should be gone after deleting the signing identity") // TODO get this right

	_ = os.RemoveAll("./testFiles")
}

func TestRemovesPublicKeyWhenRemovingSigningIdentity(t *testing.T) {
	pathToFixtures := "./testFiles"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")
	_ = manager.CreateSigningIdentity()

	manager.RemoveSigningIdentity()

	assert.Nil(t, manager.PublicKey(), "The public key should be gone after deleting the signing identity")

	_ = os.RemoveAll("./testFiles")
}

func TestRemovesPrivateKeyWhenRemovngSigningIdentity(t *testing.T) {
	pathToFixtures := "./testFiles"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt", "eckey.pem")
	_ = manager.CreateSigningIdentity()

	manager.RemoveSigningIdentity()

	assert.Nil(t, manager.PrivateKey(), "The private key should be gone after deleting the signing identity")

	_ = os.RemoveAll("./testFiles")
}
