package sentinel

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"golang.org/x/time/rate"
	. "gopkg.in/check.v1"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	ctypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

type AuthSuite struct{}

var _ = Suite(&AuthSuite{})

func (s *AuthSuite) TestArkAuth(c *C) {
	// setup
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	module.NewBasicManager().RegisterInterfaces(interfaceRegistry)
	ctypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	pubkey := types.GetRandomPubKey()
	kb := cKeys.NewInMemory(cdc)
	info, _, err := kb.NewMnemonic("whatever", cKeys.English, `m/44'/931'/0'/0/0`, "", hd.Secp256k1)
	c.Assert(err, IsNil)
	pub, err := info.GetPubKey()
	c.Assert(err, IsNil)
	pk, err := common.NewPubKeyFromCrypto(pub)
	c.Assert(err, IsNil)

	var signature []byte
	height := int64(12)
	nonce := int64(3)
	chain := common.BTCChain
	contractId := uint64(50)

	message := []byte(fmt.Sprintf("%s:%s:%s:%d:%d", pubkey.String(), chain, pk, height, nonce))
	signature, _, err = kb.Sign("whatever", message)
	c.Assert(err, IsNil)

	// happy path
	raw := GenerateArkAuthString(contractId, pk, height, nonce, signature)
	_, err = parseArkAuth(raw)
	c.Assert(err, IsNil)

	// bad signature
	raw = GenerateArkAuthString(contractId, pk, height, nonce, signature)
	_, err = parseArkAuth(raw + "randome not hex!")
	c.Assert(err, NotNil)
}

func (s *AuthSuite) TestFreeTier(c *C) {
	config := conf.Configuration{
		FreeTierRateLimitDuration: time.Minute,
		FreeTierRateLimit:         1,
	}
	proxy := NewProxy(config)

	remoteAddr := "127.0.0.1:8000"

	code, err := proxy.freeTier(remoteAddr)
	c.Assert(err, IsNil)
	c.Check(code, Equals, http.StatusOK)
	code, err = proxy.freeTier(remoteAddr)
	c.Assert(err, NotNil)
	c.Check(code, Equals, http.StatusTooManyRequests)
}

func (s *AuthSuite) TestPaidTier(c *C) {
	// setup
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	module.NewBasicManager().RegisterInterfaces(interfaceRegistry)
	ctypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	visitors = make(map[string]*rate.Limiter) // reset visitors
	pubkey := types.GetRandomPubKey()
	kb := cKeys.NewInMemory(cdc)
	info, _, err := kb.NewMnemonic("whatever", cKeys.English, `m/44'/931'/0'/0/0`, "", hd.Secp256k1)
	c.Assert(err, IsNil)
	pub, err := info.GetPubKey()
	c.Assert(err, IsNil)
	pk, err := common.NewPubKeyFromCrypto(pub)
	c.Assert(err, IsNil)

	var signature []byte
	height := int64(12)
	nonce := int64(3)
	chain := common.BTCChain.String()

	message := []byte(fmt.Sprintf("%s:%s:%s:%d:%d", pubkey.String(), chain, pk, height, nonce))
	signature, _, err = kb.Sign("whatever", message)
	c.Assert(err, IsNil)

	config := conf.Configuration{
		ProviderPubKey:            pubkey,
		AsGoTierRateLimitDuration: time.Minute,
		AsGoTierRateLimit:         1,
		SubTierRateLimitDuration:  time.Minute,
		SubTierRateLimit:          1,
		FreeTierRateLimitDuration: time.Minute,
		FreeTierRateLimit:         1,
	}
	proxy := NewProxy(config)

	contract := types.NewContract(pubkey, common.BTCChain, pk)
	contract.Height = 5
	contract.Duration = 100
	contract.Id = 545
	proxy.MemStore.SetHeight(10)
	proxy.MemStore.Put(contract)

	// happy path
	aa := ArkAuth{
		ContractId: contract.Id,
		Height:     height,
		Nonce:      nonce,
		Spender:    pk,
		Signature:  signature,
	}
	code, err := proxy.paidTier(aa, "127.0.0.1:8080")
	c.Assert(err, IsNil)
	c.Check(code, Equals, http.StatusOK)
	contract, err = proxy.MemStore.Get(contract.Key())
	c.Assert(err, IsNil)
	c.Check(contract.Nonce, Equals, int64(3))
	claim, err := proxy.ClaimStore.Get(contract.Key())
	c.Assert(err, IsNil)
	c.Check(claim.Nonce, Equals, int64(3))

	// insure that same noonce is rejected.
	code, err = proxy.paidTier(aa, "127.0.0.1:8080")
	c.Assert(err, NotNil)
	c.Check(code, Equals, http.StatusBadRequest)

	// rate limited after increasing nonce
	aa.Nonce++
	code, err = proxy.paidTier(aa, "127.0.0.1:8080")
	c.Assert(err, NotNil)
	c.Check(code, Equals, http.StatusTooManyRequests)
}

func (s *AuthSuite) TestPaidTierFailFallbackToFreeTier(c *C) {
	// setup
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	module.NewBasicManager().RegisterInterfaces(interfaceRegistry)
	ctypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	visitors = make(map[string]*rate.Limiter) // reset visitors
	pubkey := types.GetRandomPubKey()
	kb := cKeys.NewInMemory(cdc)
	info, _, err := kb.NewMnemonic("whatever", cKeys.English, `m/44'/931'/0'/0/0`, "", hd.Secp256k1)
	c.Assert(err, IsNil)
	pub, err := info.GetPubKey()
	c.Assert(err, IsNil)
	pk, err := common.NewPubKeyFromCrypto(pub)
	c.Assert(err, IsNil)

	var signature []byte
	height := int64(12)
	nonce := int64(3)
	chain := common.BTCChain.String()

	message := []byte(fmt.Sprintf("%s:%s:%s:%d:%d", pubkey.String(), chain, pk, height, nonce))
	signature, _, err = kb.Sign("whatever", message)
	c.Assert(err, IsNil)

	config := conf.Configuration{
		ProviderPubKey:            pubkey,
		AsGoTierRateLimitDuration: time.Minute,
		AsGoTierRateLimit:         1,
		SubTierRateLimitDuration:  time.Minute,
		SubTierRateLimit:          1,
		FreeTierRateLimitDuration: time.Minute,
		FreeTierRateLimit:         1,
	}
	proxy := NewProxy(config)

	contract := types.NewContract(pubkey, common.BTCChain, pk)
	contract.Height = 5
	contract.Duration = 100
	contract.Id = 55556
	// set contract to expired
	proxy.MemStore.SetHeight(120)
	proxy.MemStore.Put(contract)

	nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	})
	handlerForTest := proxy.auth(nextHandler)
	c.Assert(handlerForTest, NotNil)
	aa := ArkAuth{
		Height:     height,
		Nonce:      nonce,
		ContractId: contract.Id,
		Spender:    pk,
		Signature:  signature,
	}

	responseRecorder := httptest.NewRecorder()
	target := "/test?arkauth=" + aa.String()
	r := httptest.NewRequest(http.MethodPost, target, nil)
	handlerForTest.ServeHTTP(responseRecorder, r)
	c.Assert(responseRecorder.Code, Equals, 200)
}
