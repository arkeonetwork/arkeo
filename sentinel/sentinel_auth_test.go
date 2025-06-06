package sentinel

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

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
			assert.Equal(t, "12345", parts[0]) // contract ID
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