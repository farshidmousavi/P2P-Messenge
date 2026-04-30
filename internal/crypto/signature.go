package crypto

import (
	"crypto/ed25519"
	"errors"
)

func SignMessage(priv ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(priv, message)
}

func VerifyMessage(pub ed25519.PublicKey, message, signature []byte) error {
	if !ed25519.Verify(pub, message, signature) {
		return errors.New("invalid signature")
	}
	return nil
}
