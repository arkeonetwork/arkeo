package sentinel

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cometbft/cometbft/libs/log"
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

func TestContractAuth_Validate_TimestampUpperBound(t *testing.T) {
	client := types.GetRandomPubKey()
	lastTimestamp := int64(1000)
	
	// Create auth with timestamp too far in the future (replay attack)
	// We use time.Now() to get current time, then add more than the window
	futureTimestamp := time.Now().Unix() + ContractAuthTimestampWindow + 100 // 6+ minutes in future
	
	auth := ContractAuth{
		ContractId: 123,
		Timestamp:  futureTimestamp,
		Signature:  []byte("dummy-signature"),
		ChainId:    "test-chain",
	}
	
	// This should fail due to upper bound check (before signature validation)
	err := auth.Validate(lastTimestamp, client)
	
	// Should fail on timestamp upper bound check
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum allowed window", "Should reject future timestamp")
}

func TestContractAuth_Validate_TimestampLowerBound(t *testing.T) {
	client := types.GetRandomPubKey()
	lastTimestamp := int64(1000)
	
	// Create auth with timestamp too far in the past (replay attack)
	// This simulates a very old signature being replayed
	oldTimestamp := time.Now().Unix() - ContractAuthClockSkewTolerance - 100 // More than 1 minute in past
	
	auth := ContractAuth{
		ContractId: 123,
		Timestamp:  oldTimestamp,
		Signature:  []byte("dummy-signature"),
		ChainId:    "test-chain",
	}
	
	err := auth.Validate(lastTimestamp, client)
	
	// Should fail on timestamp lower bound check (before signature validation)
	assert.Error(t, err)
	// Could fail on either "must be larger than lastTimestamp" or "too far in the past"
	// Both are valid security checks
	if !strings.Contains(err.Error(), "must be larger than") {
		assert.Contains(t, err.Error(), "too far in the past", "Should reject very old timestamp")
	}
}

func TestContractAuth_Validate_ValidTimestampBounds(t *testing.T) {
	client := types.GetRandomPubKey()
	lastTimestamp := time.Now().Unix() - 60 // 1 minute ago
	
	// Create auth with valid timestamp (within window, after lastTimestamp)
	// Note: This will still fail signature validation, but timestamp bounds should pass
	validTimestamp := time.Now().Unix() - 30 // 30 seconds ago, within window
	
	auth := ContractAuth{
		ContractId: 123,
		Timestamp:  validTimestamp,
		Signature:  []byte("dummy-signature"),
		ChainId:    "test-chain",
	}
	
	err := auth.Validate(lastTimestamp, client)
	
	// Should fail on signature validation (expected), but NOT on timestamp bounds
	assert.Error(t, err)
	// Should NOT contain timestamp window errors - means bounds check passed
	assert.NotContains(t, err.Error(), "exceeds maximum allowed window")
	assert.NotContains(t, err.Error(), "too far in the past")
	// Should fail on signature (expected)
	assert.Contains(t, err.Error(), "invalid signature")
}

func TestContractAuth_Validate_ZeroContractId(t *testing.T) {
	client := types.GetRandomPubKey()
	
	auth := ContractAuth{
		ContractId: 0, // Invalid
		Timestamp:  1000,
		Signature:  []byte("dummy"),
		ChainId:    "test",
	}
	
	err := auth.Validate(500, client)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "contract id cannot be zero")
}

func TestIsValidIPAddress_ValidIPs(t *testing.T) {
	validIPs := []string{
		"127.0.0.1",
		"192.168.1.1",
		"10.0.0.1",
		"::1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
		"2001:db8::1",
	}
	
	for _, ip := range validIPs {
		t.Run(ip, func(t *testing.T) {
			assert.True(t, isValidIPAddress(ip), "IP %s should be valid", ip)
		})
	}
}

func TestIsValidIPAddress_InvalidIPs(t *testing.T) {
	invalidIPs := []string{
		"",
		"not.an.ip",
		"256.256.256.256",
		"192.168.1",
		"192.168.1.1.1",
		"invalid",
		"localhost", // Not a valid IP format (though we check for it separately)
	}
	
	for _, ip := range invalidIPs {
		t.Run(ip, func(t *testing.T) {
			if ip == "localhost" {
				// localhost is not a valid IP but we handle it separately in isLocalhost
				t.Skip("localhost handled separately")
			}
			assert.False(t, isValidIPAddress(ip), "IP %s should be invalid", ip)
		})
	}
}

func TestIsTrustedProxy_DefaultLocalhost(t *testing.T) {
	// Test default behavior (no trusted IPs configured)
	trustedIPs := []string{} // Empty = default to localhost
	
	tests := []struct {
		name       string
		remoteAddr string
		expected   bool
	}{
		{"localhost IPv4", "127.0.0.1:12345", true},
		{"localhost IPv6", "[::1]:12345", true},
		{"localhost string", "localhost:12345", true},
		{"external IP", "192.168.1.1:12345", false},
		{"public IP", "8.8.8.8:12345", false},
		{"empty", "", false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTrustedProxy(tt.remoteAddr, trustedIPs)
			assert.Equal(t, tt.expected, result, "isTrustedProxy(%q, %v) = %v, want %v", 
				tt.remoteAddr, trustedIPs, result, tt.expected)
		})
	}
}

func TestIsTrustedProxy_ConfiguredIPs(t *testing.T) {
	// Test with configured trusted proxy IPs
	trustedIPs := []string{"10.0.0.1", "192.168.0.0/16", "172.16.0.1"}
	
	tests := []struct {
		name       string
		remoteAddr string
		expected   bool
	}{
		{"exact match", "10.0.0.1:12345", true},
		{"CIDR match", "192.168.1.1:12345", true},
		{"CIDR match boundary", "192.168.255.255:12345", true},
		{"CIDR no match", "192.169.1.1:12345", false},
		{"exact match 172", "172.16.0.1:12345", true},
		{"no match", "8.8.8.8:12345", false},
		{"localhost still works", "127.0.0.1:12345", false}, // Not in config, so false
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isTrustedProxy(tt.remoteAddr, trustedIPs)
			assert.Equal(t, tt.expected, result, "isTrustedProxy(%q, %v) = %v, want %v", 
				tt.remoteAddr, trustedIPs, result, tt.expected)
		})
	}
}

func TestGetRemoteAddr_DirectConnection(t *testing.T) {
	// Test direct connection (no proxy)
	proxy := Proxy{
		Config: conf.Configuration{
			TrustedProxyIPs: []string{}, // Empty = localhost-only
		},
	}
	
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "203.0.113.1:54321" // External IP
	
	// Attacker tries to spoof X-Forwarded-For
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	req.Header.Set("X-Real-Ip", "1.2.3.4")
	
	ip := proxy.getRemoteAddr(req)
	
	// Should ignore spoofed headers and use RemoteAddr
	assert.Equal(t, "203.0.113.1", ip, "Should use RemoteAddr, not spoofed header")
}

func TestGetRemoteAddr_BehindTrustedProxy(t *testing.T) {
	// Test behind trusted proxy
	proxy := Proxy{
		Config: conf.Configuration{
			TrustedProxyIPs: []string{"127.0.0.1"}, // Trust localhost
		},
	}
	
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:54321" // From localhost proxy
	
	// Proxy sets real client IP
	req.Header.Set("X-Real-Ip", "192.168.1.100")
	
	ip := proxy.getRemoteAddr(req)
	
	// Should trust header from trusted proxy
	assert.Equal(t, "192.168.1.100", ip, "Should use X-Real-Ip from trusted proxy")
}

func TestGetRemoteAddr_XForwardedFor_MultipleIPs(t *testing.T) {
	// Test X-Forwarded-For with multiple IPs (client, proxy1, proxy2)
	proxy := Proxy{
		Config: conf.Configuration{
			TrustedProxyIPs: []string{"10.0.0.1"},
		},
	}
	
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "10.0.0.1:54321" // From trusted proxy
	
	// X-Forwarded-For format: "client, proxy1, proxy2"
	req.Header.Set("X-Forwarded-For", "203.0.113.50, 10.0.0.2, 10.0.0.1")
	
	ip := proxy.getRemoteAddr(req)
	
	// Should take first IP (original client)
	assert.Equal(t, "203.0.113.50", ip, "Should extract first IP from X-Forwarded-For")
}

func TestGetRemoteAddr_InvalidIPInHeader(t *testing.T) {
	// Test with invalid IP in header
	proxy := Proxy{
		Config: conf.Configuration{
			TrustedProxyIPs: []string{"127.0.0.1"},
		},
	}
	
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "127.0.0.1:54321"
	
	// Invalid IP format
	req.Header.Set("X-Real-Ip", "not.an.ip.address")
	
	ip := proxy.getRemoteAddr(req)
	
	// Should fallback to RemoteAddr
	assert.Equal(t, "127.0.0.1", ip, "Should fallback to RemoteAddr on invalid header IP")
}

func TestGetRemoteAddr_SpoofingPrevention(t *testing.T) {
	// Test that spoofed headers from untrusted source are ignored
	proxy := Proxy{
		Config: conf.Configuration{
			TrustedProxyIPs: []string{"10.0.0.1"}, // Only trust this IP
		},
	}
	
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "203.0.113.1:54321" // Attacker's real IP (not trusted)
	
	// Attacker tries to spoof
	req.Header.Set("X-Forwarded-For", "1.2.3.4") // Whitelisted IP
	req.Header.Set("X-Real-Ip", "1.2.3.4")
	
	ip := proxy.getRemoteAddr(req)
	
	// Should ignore headers and use RemoteAddr
	assert.Equal(t, "203.0.113.1", ip, "Should ignore spoofed headers from untrusted source")
}
