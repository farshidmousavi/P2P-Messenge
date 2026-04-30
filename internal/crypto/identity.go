package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"os"
)

type Identity struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func GenerateIdentity() (*Identity, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	return &Identity{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

func (i *Identity) ID() string {
	return hex.EncodeToString(i.PublicKey)
}

func (i *Identity) Save(path string) error {
	if _, err := os.Stat(path); err == nil {
		return errors.New("identity already exists")
	}

	data := append(i.PrivateKey, i.PublicKey...)
	return os.WriteFile(path, data, 0600)
}

func LoadIdentity(path string) (*Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if len(data) != ed25519.PrivateKeySize+ed25519.PublicKeySize {
		return nil, errors.New("invalid identity file")
	}

	priv := ed25519.PrivateKey(data[:ed25519.PrivateKeySize])
	pub := ed25519.PublicKey(data[ed25519.PrivateKeySize:])

	return &Identity{
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}
