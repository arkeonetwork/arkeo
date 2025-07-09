package sentinel

import (
	"fmt"
	"sync"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
)

type ArkeoAuthManager struct {
	contractId uint64
	chainId    string
	privKey    *secp256k1.PrivKey
	nonce      int64
	nonceStore *NonceStore
	mu         sync.Mutex
	logger     log.Logger
}

func NewArkeoAuthManager(contractId uint64, chainId string, mnemonic string, nonceStore *NonceStore, logger log.Logger) (*ArkeoAuthManager, error) {
	if contractId == 0 || chainId == "" || mnemonic == "" {
		return nil, nil // Auth not configured
	}

	// Get default HD path (same as in signThis function)
	hdPath := hd.NewFundraiserParams(0, 118, 0).String()

	// Derive private key from mnemonic
	derivedPriv, err := hd.Secp256k1.Derive()(mnemonic, "", hdPath)
	if err != nil {
		return nil, fmt.Errorf("failed to derive private key: %w", err)
	}

	privKey := hd.Secp256k1.Generate()(derivedPriv).(*secp256k1.PrivKey)

	// Load last nonce from store
	lastNonce := int64(0)
	if nonceStore != nil {
		lastNonce, err = nonceStore.Get(contractId)
		if err != nil {
			logger.Error("failed to load nonce from store", "error", err)
			// Continue with nonce 0 if load fails
		} else {
			logger.Info("loaded nonce from store", "contractId", contractId, "nonce", lastNonce)
		}
	}

	return &ArkeoAuthManager{
		contractId: contractId,
		chainId:    chainId,
		privKey:    privKey,
		nonce:      lastNonce,
		nonceStore: nonceStore,
		logger:     logger,
	}, nil
}

func (am *ArkeoAuthManager) GenerateAuthHeader() (string, error) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.nonce++

	// Persist new nonce
	if am.nonceStore != nil {
		if err := am.nonceStore.Set(am.contractId, am.nonce); err != nil {
			am.logger.Error("failed to persist nonce", "error", err)
			// Continue even if persistence fails
		}
	}

	// Generate message to sign (using existing function from sentinel_auth.go)
	message := GenerateMessageToSign(am.contractId, am.nonce, am.chainId)

	// Sign the message
	sig, err := am.privKey.Sign([]byte(message))
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %w", err)
	}

	// Generate auth string (using existing function from sentinel_auth.go)
	authString := GenerateArkAuthString(am.contractId, am.nonce, sig, am.chainId)

	am.logger.Debug("generated auth header", "contractId", am.contractId, "nonce", am.nonce)

	return authString, nil
}

func (am *ArkeoAuthManager) GetNonce() int64 {
	am.mu.Lock()
	defer am.mu.Unlock()
	return am.nonce
}

func (am *ArkeoAuthManager) GetContractId() uint64 {
	return am.contractId
}

func (am *ArkeoAuthManager) GetChainId() string {
	return am.chainId
}

func (am *ArkeoAuthManager) Close() error {
	if am.nonceStore != nil {
		return am.nonceStore.Close()
	}
	return nil
}

// GetPublicKey returns the public key derived from the private key
func (am *ArkeoAuthManager) GetPublicKey() secp256k1.PubKey {
	return *am.privKey.PubKey().(*secp256k1.PubKey)
}
