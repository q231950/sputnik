package keymanager

import (
	"crypto/ecdsa"
)

type MockKeyManager struct {
}

func (m MockKeyManager) PublicKey() *ecdsa.PublicKey {
	return nil
}

func (m MockKeyManager) PrivateKey() *ecdsa.PrivateKey {
	return nil
}

func (m MockKeyManager) KeyId() string {
	return "key id"
}
