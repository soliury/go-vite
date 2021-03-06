package keystore

import (
	"errors"
	"github.com/pborman/uuid"
	"github.com/vitelabs/go-vite/common/types"
	"github.com/vitelabs/go-vite/crypto/ed25519"
	"strconv"
)

const (
	keystoreVersion = 1
)

type keyStore interface {
	// Returns the key associated with the given address , using the given password to recover it from a file.
	ExtractKey(address types.Address, password string) (*Key, error)

	StoreKey(k *Key, password string) error
}

type Key struct {
	Id         uuid.UUID
	Address    types.Address
	PublicKey  *ed25519.PublicKey
	PrivateKey *ed25519.PrivateKey
}

type encryptedKeyJSON struct {
	HexAddress string     `json:"hexaddress"`
	HexPubKey  string     `json:"hexpubkey"`
	Id         string     `json:"id"`
	Crypto     cryptoJSON `json:"crypto"`
	Version    int        `json:"keystoreversion"`
	Timestamp  int64      `json:"timestamp"`
}

type cryptoJSON struct {
	CipherName   string       `json:"ciphername"`
	CipherText   string       `json:"ciphertext"`
	Nonce        string       `json:"nonce"`
	KDF          string       `json:"kdf"`
	ScryptParams scryptParams `json:"scryptparams"`
}

type scryptParams struct {
	N      int    `json:"n"`
	R      int    `json:"r"`
	P      int    `json:"p"`
	KeyLen int    `json:"keylen"`
	Salt   string `json:"salt"`
}

func (key *Key) Sign(data []byte) ([]byte, error) {
	if l := len(*key.PrivateKey); l != ed25519.PrivateKeySize {
		return nil, errors.New("ed25519: bad private key length: " + strconv.Itoa(l))
	}
	return ed25519.Sign(*key.PrivateKey, data), nil
}
