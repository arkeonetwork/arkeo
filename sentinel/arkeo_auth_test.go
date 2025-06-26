package sentinel

import (
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testMnemonic = strings.Repeat("dog ", 23) + "fossil" // Same format as regression tests
	testChainId  = "arkeo-testnet"
)

func TestNewArkeoAuthManager_NoConfig(t *testing.T) {
	logger := log.NewNopLogger()

	// Test with zero contract ID
	am, err := NewArkeoAuthManager(0, testChainId, testMnemonic, nil, logger)
	assert.NoError(t, err)
	assert.Nil(t, am)

	// Test with empty chain ID
	am, err = NewArkeoAuthManager(12345, "", testMnemonic, nil, logger)
	assert.NoError(t, err)
	assert.Nil(t, am)

	// Test with empty mnemonic
	am, err = NewArkeoAuthManager(12345, testChainId, "", nil, logger)
	assert.NoError(t, err)
	assert.Nil(t, am)
}

func TestNewArkeoAuthManager_InvalidMnemonic(t *testing.T) {
	logger := log.NewNopLogger()

	// Test with invalid mnemonic
	_, err := NewArkeoAuthManager(12345, testChainId, "invalid mnemonic phrase", nil, logger)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to derive private key")
}

func TestNewArkeoAuthManager_ValidConfig(t *testing.T) {
	logger := log.NewNopLogger()
	store, err := NewNonceStore("")
	require.NoError(t, err)
	defer store.Close()

	am, err := NewArkeoAuthManager(12345, testChainId, testMnemonic, store, logger)
	assert.NoError(t, err)
	assert.NotNil(t, am)
	assert.Equal(t, uint64(12345), am.GetContractId())
	assert.Equal(t, testChainId, am.GetChainId())
	assert.Equal(t, int64(0), am.GetNonce())
}

func TestArkeoAuthManager_GenerateAuthHeader(t *testing.T) {
	logger := log.NewNopLogger()
	store, err := NewNonceStore("")
	require.NoError(t, err)
	defer store.Close()

	am, err := NewArkeoAuthManager(12345, testChainId, testMnemonic, store, logger)
	require.NoError(t, err)

	// Generate first auth header
	authHeader1, err := am.GenerateAuthHeader()
	assert.NoError(t, err)
	assert.NotEmpty(t, authHeader1)

	// Verify format: contractId:nonce:chainId:signature
	parts := strings.Split(authHeader1, ":")
	assert.Len(t, parts, 4)
	assert.Equal(t, "12345", parts[0])
	assert.Equal(t, "1", parts[1])
	assert.Equal(t, testChainId, parts[2])
	assert.NotEmpty(t, parts[3]) // signature

	// Generate second auth header - nonce should increment
	authHeader2, err := am.GenerateAuthHeader()
	assert.NoError(t, err)
	assert.NotEqual(t, authHeader1, authHeader2)

	parts2 := strings.Split(authHeader2, ":")
	assert.Equal(t, "2", parts2[1])
}

func TestArkeoAuthManager_NoncePersistence(t *testing.T) {
	logger := log.NewNopLogger()
	tmpDir, err := os.MkdirTemp("", "auth-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create auth manager and generate some headers
	store1, err := NewNonceStore(tmpDir)
	require.NoError(t, err)

	am1, err := NewArkeoAuthManager(12345, testChainId, testMnemonic, store1, logger)
	require.NoError(t, err)

	// Generate 5 auth headers
	for i := 0; i < 5; i++ {
		_, err = am1.GenerateAuthHeader()
		assert.NoError(t, err)
	}

	assert.Equal(t, int64(5), am1.GetNonce())

	// Close and reopen
	err = am1.Close()
	assert.NoError(t, err)

	// Create new auth manager with same store location
	store2, err := NewNonceStore(tmpDir)
	require.NoError(t, err)

	am2, err := NewArkeoAuthManager(12345, testChainId, testMnemonic, store2, logger)
	require.NoError(t, err)
	defer am2.Close()

	// Should start from last saved nonce
	assert.Equal(t, int64(5), am2.GetNonce())

	// Generate next header should be nonce 6
	authHeader, err := am2.GenerateAuthHeader()
	assert.NoError(t, err)

	parts := strings.Split(authHeader, ":")
	assert.Equal(t, "6", parts[1])
}

func TestArkeoAuthManager_ConcurrentGeneration(t *testing.T) {
	logger := log.NewNopLogger()
	store, err := NewNonceStore("")
	require.NoError(t, err)
	defer store.Close()

	am, err := NewArkeoAuthManager(12345, testChainId, testMnemonic, store, logger)
	require.NoError(t, err)

	// Generate headers concurrently
	numGoroutines := 100
	var wg sync.WaitGroup
	nonces := make(map[string]bool)
	var mu sync.Mutex

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			authHeader, err := am.GenerateAuthHeader()
			assert.NoError(t, err)

			parts := strings.Split(authHeader, ":")
			mu.Lock()
			nonces[parts[1]] = true
			mu.Unlock()
		}()
	}

	wg.Wait()

	// All nonces should be unique
	assert.Len(t, nonces, numGoroutines)
	assert.Equal(t, int64(numGoroutines), am.GetNonce())
}

func TestArkeoAuthManager_PublicKey(t *testing.T) {
	logger := log.NewNopLogger()
	
	am, err := NewArkeoAuthManager(12345, testChainId, testMnemonic, nil, logger)
	require.NoError(t, err)

	pubKey := am.GetPublicKey()
	assert.NotNil(t, pubKey)
	
	// Verify we can generate consistent signatures
	authHeader1, err := am.GenerateAuthHeader()
	assert.NoError(t, err)
	
	// Extract signature and verify format
	parts := strings.Split(authHeader1, ":")
	assert.Len(t, parts, 4)
}

func TestArkeoAuthManager_MultipleContracts(t *testing.T) {
	logger := log.NewNopLogger()
	store, err := NewNonceStore("")
	require.NoError(t, err)
	defer store.Close()

	// Create two auth managers for different contracts
	am1, err := NewArkeoAuthManager(12345, testChainId, testMnemonic, store, logger)
	require.NoError(t, err)

	am2, err := NewArkeoAuthManager(67890, testChainId, testMnemonic, store, logger)
	require.NoError(t, err)

	// Generate headers for both
	_, err = am1.GenerateAuthHeader()
	assert.NoError(t, err)

	_, err = am2.GenerateAuthHeader()
	assert.NoError(t, err)

	// Nonces should be independent
	assert.Equal(t, int64(1), am1.GetNonce())
	assert.Equal(t, int64(1), am2.GetNonce())

	// Generate more for contract 1
	_, err = am1.GenerateAuthHeader()
	assert.NoError(t, err)
	_, err = am1.GenerateAuthHeader()
	assert.NoError(t, err)

	assert.Equal(t, int64(3), am1.GetNonce())
	assert.Equal(t, int64(1), am2.GetNonce())
}