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
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt")
	assert.Equal(t, "abc\n", manager.KeyID(), "A stored key id should get found")
}

func TestPrivateKey(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt")

	expectedD := new(big.Int)
	expectedD.SetString("57359333433306843951573675484597381433848383364258304847182053853963006392866", 10)
	assert.Equal(t, expectedD, manager.PrivateKey().D, "The private key does not match the certificate")
}

func TestPublicKey(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt")

	expectedX := new(big.Int)
	expectedX.SetString("83764337057786748884235593670888306280068598385338703097790983192099570881279", 10)
	assert.Equal(t, expectedX, manager.PublicKey().X, "The public key's X does not match the certificate")

	expectedY := new(big.Int)
	expectedY.SetString("52205868503601725522780045291819081986668283346034058232083577220213552372588", 10)
	assert.Equal(t, expectedY, manager.PublicKey().Y, "The public key's Y does not match the certificate")
}

func TestStoreKeyID(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyidTest.txt")

	manager.StoreKeyID("key")
	assert.Equal(t, "key", manager.KeyID(), "The key id should be stored correctly")

	_ = os.Remove(pathToFixtures + "/keyidTest.txt")
}

func TestSigningIdentityExists(t *testing.T) {
	pathToFixtures := "./fixtures"
	manager := NewWithSecretsFolder(pathToFixtures, "keyid.txt")
	assert.True(t, manager.SigningIdentityExists(), "The signing identity should exist in this fixture")
}
