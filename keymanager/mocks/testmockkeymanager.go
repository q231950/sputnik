package keymanager

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type MockKeyManager struct {
}

func (m MockKeyManager) PublicKey() *ecdsa.PublicKey {
	return nil
}

func (m MockKeyManager) PrivateKey() *ecdsa.PrivateKey {
	c := elliptic.P256()
	key, _ := ecdsa.GenerateKey(c, rand.Reader)
	return key
}

func (m MockKeyManager) KeyId() string {
	return "key id"
}

func (m MockKeyManager) RemoveSigningIdentity() error {
	return nil
}

func (m MockKeyManager) StoreKeyId(key string) error {
	return nil
}
