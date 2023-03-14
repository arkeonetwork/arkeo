package sentinel

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

type MemStoreSuite struct {
	suite.Suite
	server *httptest.Server
}

func (s *MemStoreSuite) SetUpTest() {
	s.server = httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch {
		case strings.HasSuffix(req.RequestURI, "/arkeo/contract/cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu/arkeo-mainnet/cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu"):
			httpTestHandler(s.T(), rw, `
{ "contract": {
				"provider_pub_key": "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu",
				"service": 1,
				"client": "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu",
				"delegate": "cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu",
				"type": 0,
				"height": "15",
				"duration": "100",
				"rate": "3",
				"deposit": "500",
				"paid": "0",
				"nonce": "9",
				"settlement_height": "0"
			}}`)
		default:
			fmt.Println(req.RequestURI)
			panic("could not serve request")
		}
	}))
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
	var err error
	baseURL := fmt.Sprintf("http://%s", s.server.Listener.Addr().String())
	mem := NewMemStore(baseURL, log.NewTMLogger(log.NewSyncWriter(os.Stdout)))

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
	pk, err := common.NewPubKey("cosmospub1addwnpepqg3523h7e7ggeh6na2lsde6s394tqxnvufsz0urld6zwl8687ue9c3dasgu")
	require.NoError(s.T(), err)

	key := mem.Key(pk.String(), "arkeo-mainnet", pk.String())
	contract, err = mem.Get(key)
	require.NoError(s.T(), err)
	require.Equal(s.T(), contract.Rate, int64(3))
	require.Equal(s.T(), contract.Deposit.Int64(), int64(500))
	require.Equal(s.T(), contract.Paid.Int64(), int64(0))
}
