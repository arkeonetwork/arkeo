package sentinel

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

type MemStoreSuite struct {
	suite.Suite
	server *httptest.Server
}

// SetUpTest is removed - tests now create their own servers

func (s *MemStoreSuite) TearDownTest() {
	if s.server != nil {
		s.server.Close()
	}
}

func httpTestHandler(t *testing.T, rw http.ResponseWriter, content string) {
	if content == "500" {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	if _, err := rw.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
}

func (s *MemStoreSuite) TestMemStore() {
	// Use a dynamic test pubkey to avoid hardcoding invalid ones
	testPK := types.GetRandomPubKey()
	
	// Recreate server with the dynamic pubkey
	s.server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		expectedURI := fmt.Sprintf("/arkeo/contract/%s/arkeo-mainnet/%s", testPK.String(), testPK.String())
		switch {
		case strings.HasSuffix(req.RequestURI, expectedURI):
			httpTestHandler(s.T(), rw, fmt.Sprintf(`
{ "contract": {
				"provider_pub_key": "%s",
				"service": 1,
				"client": "%s",
				"delegate": "%s",
				"type": 0,
				"height": "15",
				"duration": "100",
				"rate": {"denom": "uarkeo", "amount": "3"},
				"deposit": "500",
				"paid": "0",
				"nonce": "9",
				"settlement_height": "0"
			}}`, testPK.String(), testPK.String(), testPK.String()))
		default:
			panic(fmt.Sprintf("could not serve request: %s", req.RequestURI))
		}
	}))
	defer s.server.Close()
	
	var err error
	baseURL := fmt.Sprintf("http://%s", s.server.Listener.Addr().String())
	mem := NewMemStore(baseURL, nil, log.NewTMLogger(log.NewSyncWriter(os.Stdout)))

	require.Equal(s.T(), mem.Key("foo", "bar", "baz"), "foo/bar/baz")

	mem.SetHeight(30)
	require.Equal(s.T(), mem.GetHeight(), int64(30))

	pk1 := types.GetRandomPubKey()
	pk2 := types.GetRandomPubKey()
	contract := types.NewContract(pk1, common.Service(0), pk2)
	contract.Height = 4
	contract.Duration = 100
	contract.Id = 55786

	mem.Put(contract)

	contract, err = mem.Get(contract.Key())
	require.NoError(s.T(), err)
	require.Equal(s.T(), contract.Height, int64(4))

	// fetch from server
	key := mem.Key(testPK.String(), "arkeo-mainnet", testPK.String())
	contract, err = mem.Get(key)
	require.NoError(s.T(), err)
	require.Equal(s.T(), contract.Rate.Amount.Int64(), int64(3))
	require.Equal(s.T(), contract.Deposit.Int64(), int64(500))
	require.Equal(s.T(), contract.Paid.Int64(), int64(0))
}

func (s *MemStoreSuite) TestMemStoreWithAuth() {
	// Use a dynamic test pubkey
	testPK := types.GetRandomPubKey()
	
	// Create a test server that verifies auth header
	authChecked := false
	testServer := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get(QueryArkAuth)
		if authHeader != "" {
			authChecked = true
			// Verify auth header format
			parts := strings.Split(authHeader, ":")
			require.Len(s.T(), parts, 4)
			require.Equal(s.T(), "12345", parts[0]) // contract ID
			require.Equal(s.T(), "1", parts[1])     // nonce
			require.Equal(s.T(), "test-chain", parts[2]) // chain ID
		}
		
		// Return a mock contract
		httpTestHandler(s.T(), rw, fmt.Sprintf(`
{ "contract": {
				"provider_pub_key": "%s",
				"service": 1,
				"client": "%s",
				"delegate": "%s",
				"type": 0,
				"height": "15",
				"duration": "100",
				"rate": {"denom": "uarkeo", "amount": "3"},
				"deposit": "500",
				"paid": "0",
				"nonce": "9",
				"settlement_height": "0"
			}}`, testPK.String(), testPK.String(), testPK.String()))
	}))
	defer testServer.Close()

	// Create auth manager
	logger := log.NewNopLogger()
	nonceStore, err := NewNonceStore("")
	require.NoError(s.T(), err)
	defer nonceStore.Close()
	
	testMnemonic := strings.Repeat("dog ", 23) + "fossil"
	authManager, err := NewArkeoAuthManager(12345, "test-chain", testMnemonic, nonceStore, logger)
	require.NoError(s.T(), err)

	// Create memstore with auth
	mem := NewMemStore(testServer.URL, authManager, logger)

	// Fetch contract - should add auth header
	key := mem.Key(testPK.String(), "arkeo-mainnet", testPK.String())
	contract, err := mem.Get(key)
	require.NoError(s.T(), err)
	require.Equal(s.T(), contract.Rate.Amount.Int64(), int64(3))
	
	// Verify auth header was checked
	require.True(s.T(), authChecked, "Auth header should have been sent")
}

func TestMemStoreSuite(t *testing.T) {
	suite.Run(t, new(MemStoreSuite))
}
