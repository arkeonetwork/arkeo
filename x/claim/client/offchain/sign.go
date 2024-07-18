package offchain

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

func signData(data string, privateKey string) ([]byte, error) {

	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return nil, err
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	hash := sha256.Sum256([]byte(data))
	signature, err := privKey.Sign(hash[:])
	if err != nil {
		return nil, err
	}

	return signature, nil
}
