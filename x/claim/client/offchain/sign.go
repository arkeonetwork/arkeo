package offchain

import (
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

type SignedData struct {
	Signature   string `json:"proof_signature"`
	ProofPubkey string `json:"proof_pubkey"`
}

func (s *SignedData) sign(data string, privateKey string) error {
	privKeyBytes, err := hex.DecodeString(privateKey)
	if err != nil {
		return err
	}

	privKey := secp256k1.PrivKey{Key: privKeyBytes}
	hash := sha512.Sum512([]byte(data))
	signature, err := privKey.Sign(hash[:])
	if err != nil {
		return err
	}

	s.Signature = hex.EncodeToString(signature)

	const prefix = "PubKeySecp256k1{"
	const suffix = "}"

	if !strings.HasPrefix(privKey.PubKey().String(), prefix) || !strings.HasSuffix(privKey.PubKey().String(), suffix) {
		return fmt.Errorf("invalid pubkey format")
	}

	publicKey := strings.TrimPrefix(privKey.PubKey().String(), prefix)
	publicKey = strings.TrimSuffix(publicKey, suffix)

	s.ProofPubkey = publicKey
	return nil
}

func (s *SignedData) getSignedDataString() (string, error) {

	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshall signed data: %w", err)
	}

	return string(jsonData), nil
}
