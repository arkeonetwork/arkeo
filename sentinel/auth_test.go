package sentinel

import (
	"encoding/hex"
	"fmt"
	"net/http"
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

type AuthSuite struct {
}

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
	chain := common.BTCChain.String()

	message := []byte(fmt.Sprintf("%s:%s:%s:%d:%d", pubkey.String(), chain, pk, height, nonce))
	signature, _, err = kb.Sign("whatever", message)
	c.Assert(err, IsNil)

	sig := hex.EncodeToString(signature)

	// happy path
	raw := fmt.Sprintf("%s:%s:%s:%d:%d:%s", pubkey.String(), chain, pk, height, nonce, sig)
	_, err = parseArkAuth(raw)
	c.Assert(err, IsNil)

	// bad signature
	raw = fmt.Sprintf("%s:%s:%s:%d:%d:%s", pubkey.String(), chain, pk, height, nonce, "bad siggy")
	_, err = parseArkAuth(raw)
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
	proxy.MemStore.SetHeight(10)
	key := proxy.MemStore.Key(pubkey.String(), common.BTCChain.String(), pk.String())
	proxy.MemStore.Put(key, contract)

	// happy path
	aa := ArkAuth{
		Provider:  pubkey,
		Height:    height,
		Nonce:     nonce,
		Chain:     common.BTCChain,
		Spender:   pk,
		Signature: signature,
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

	// rate limited
	code, err = proxy.paidTier(aa, "127.0.0.1:8080")
	c.Assert(err, NotNil)
	c.Check(code, Equals, http.StatusTooManyRequests)
}
