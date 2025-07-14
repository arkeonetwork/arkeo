package sentinel

import (
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cometbft/cometbft/libs/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// If both files define testMnemonic, keep only one definition.
// const testMnemonic = "swallow century piece actor harsh foam remove web range harsh frost powder flock category subway release object brave syrup ginger echo west phrase voice"

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

func TestHandleActiveContract_WithAuth(t *testing.T) {
	// Create a test server that verifies auth header
	authChecked := false
	var receivedAuthHeader string
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedAuthHeader = r.Header.Get(QueryArkAuth)
		if receivedAuthHeader != "" {
			authChecked = true
			// Verify auth header format
			parts := strings.Split(receivedAuthHeader, ":")
			assert.Len(t, parts, 4)
			assert.Equal(t, "12345", parts[0])      // contract ID
			assert.Equal(t, "arkeo-test", parts[2]) // chain ID
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"contract": {"id": "1"}}`))
	}))
	defer testServer.Close()

	// Create auth manager
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	nonceStore, err := NewNonceStore("")
	require.NoError(t, err)
	defer nonceStore.Close()

	authManager, err := NewArkeoAuthManager(12345, "arkeo-test", testMnemonic, nonceStore, logger)
	require.NoError(t, err)

	// Create proxy with auth
	config := conf.Configuration{
		ProviderPubKey:      types.GetRandomPubKey(),
		ArkeoAuthContractId: 12345,
		ArkeoAuthChainId:    "arkeo-test",
	}

	proxy := Proxy{
		Config:      config,
		logger:      logger,
		authManager: authManager,
		proxies: map[string]*url.URL{
			"arkeo-mainnet-fullnode": common.MustParseURL(testServer.URL),
		},
	}

	// Create test request
	req := httptest.NewRequest("GET", "/arkeo/active-contract", nil)
	req = mux.SetURLVars(req, map[string]string{
		"service": "btc-mainnet-fullnode",
		"spender": "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu",
	})

	w := httptest.NewRecorder()
	proxy.handleActiveContract(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify auth header was sent
	assert.True(t, authChecked, "Auth header should have been sent")
	assert.NotEmpty(t, receivedAuthHeader)
}

func TestCreateAuthenticatedReverseProxy(t *testing.T) {
	// Create a test backend server
	authHeaders := []string{}
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get(QueryArkAuth); auth != "" {
			authHeaders = append(authHeaders, auth)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))
	defer backend.Close()

	// Create auth manager
	logger := log.NewNopLogger()
	nonceStore, err := NewNonceStore("")
	require.NoError(t, err)
	defer nonceStore.Close()

	authManager, err := NewArkeoAuthManager(12345, "arkeo-test", testMnemonic, nonceStore, logger)
	require.NoError(t, err)

	// Create proxy
	proxy := Proxy{
		logger:      logger,
		authManager: authManager,
	}

	// Create authenticated reverse proxy
	target := common.MustParseURL(backend.URL)
	reverseProxy := proxy.createAuthenticatedReverseProxy(target)

	// Create test server using the reverse proxy
	testServer := httptest.NewServer(reverseProxy)
	defer testServer.Close()

	// Make multiple requests
	for i := 0; i < 3; i++ {
		resp, err := http.Get(testServer.URL)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Verify auth headers were sent
	assert.Len(t, authHeaders, 3)

	// Verify nonces increment
	for i, authHeader := range authHeaders {
		parts := strings.Split(authHeader, ":")
		assert.Equal(t, strconv.Itoa(i+1), parts[1]) // nonce should be 1, 2, 3
	}
}

func TestProxy_NoAuth(t *testing.T) {
	// Create a test backend that should NOT receive auth headers
	backend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get(QueryArkAuth), "Should not have auth header")
		w.WriteHeader(http.StatusOK)
	}))
	defer backend.Close()

	// Create proxy without auth manager
	logger := log.NewNopLogger()
	proxy := Proxy{
		logger:      logger,
		authManager: nil, // No auth configured
	}

	// Create reverse proxy
	target := common.MustParseURL(backend.URL)
	reverseProxy := proxy.createAuthenticatedReverseProxy(target)

	// Create test server
	testServer := httptest.NewServer(reverseProxy)
	defer testServer.Close()

	// Make request
	resp, err := http.Get(testServer.URL)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}
