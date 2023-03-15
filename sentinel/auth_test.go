package sentinel

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"

	"golang.org/x/time/rate"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	ctypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

func TestArkAuth(t *testing.T) {
	// setup
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	module.NewBasicManager().RegisterInterfaces(interfaceRegistry)
	ctypes.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	pubkey := types.GetRandomPubKey()
	kb := cKeys.NewInMemory(cdc)
	info, _, err := kb.NewMnemonic("whatever", cKeys.English, `m/44'/931'/0'/0/0`, "", hd.Secp256k1)
	require.NoError(t, err)
	pub, err := info.GetPubKey()
	require.NoError(t, err)
	pk, err := common.NewPubKeyFromCrypto(pub)
	require.NoError(t, err)

	var signature []byte
	nonce := int64(3)
	service := common.BTCService
	contractId := uint64(50)

	message := []byte(fmt.Sprintf("%s:%s:%s:%d", pubkey.String(), service, pk, nonce))
	signature, _, err = kb.Sign("whatever", message)
	require.NoError(t, err)

	// happy path
	raw := GenerateArkAuthString(contractId, pk, nonce, signature)
	_, err = parseArkAuth(raw)
	require.NoError(t, err)

	// bad signature
	raw = GenerateArkAuthString(contractId, pk, nonce, signature)
	_, err = parseArkAuth(raw + "randome not hex!")
	require.Error(t, err)
}

func TestFreeTier(t *testing.T) {
	config := conf.Configuration{
		FreeTierRateLimitDuration: time.Minute,
		FreeTierRateLimit:         1,
	}
	proxy := NewProxy(config)

	remoteAddr := "127.0.0.1:8000"

	code, err := proxy.freeTier(remoteAddr)
	require.NoError(t, err)
	require.Equal(t, code, http.StatusOK)

	code, err = proxy.freeTier(remoteAddr)
	require.Error(t, err)
	require.Equal(t, code, http.StatusTooManyRequests)
}

func TestPaidTier(t *testing.T) {
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
	require.NoError(t, err)
	pub, err := info.GetPubKey()
	require.NoError(t, err)
	pk, err := common.NewPubKeyFromCrypto(pub)
	require.NoError(t, err)

	var signature []byte
	nonce := int64(3)
	service := common.BTCService.String()

	message := []byte(fmt.Sprintf("%s:%s:%s:%d", pubkey.String(), service, pk, nonce))
	signature, _, err = kb.Sign("whatever", message)
	require.NoError(t, err)

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

	contract := types.NewContract(pubkey, common.BTCService, pk)
	contract.Height = 5
	contract.Duration = 100
	contract.Id = 545
	proxy.MemStore.SetHeight(10)
	proxy.MemStore.Put(contract)

	// happy path
	aa := ArkAuth{
		ContractId: contract.Id,
		Nonce:      nonce,
		Spender:    pk,
		Signature:  signature,
	}
	code, err := proxy.paidTier(aa, "127.0.0.1:8080")
	require.NoError(t, err)
	require.Equal(t, code, http.StatusOK)
	contract, err = proxy.MemStore.Get(contract.Key())
	require.NoError(t, err)
	require.Equal(t, contract.Nonce, int64(3))
	claim, err := proxy.ClaimStore.Get(contract.Key())
	require.NoError(t, err)
	require.Equal(t, claim.Nonce, int64(3))

	// insure that same noonce is rejected.
	code, err = proxy.paidTier(aa, "127.0.0.1:8080")
	require.Error(t, err)
	require.Equal(t, code, http.StatusBadRequest)

	// rate limited after increasing nonce
	aa.Nonce++
	code, err = proxy.paidTier(aa, "127.0.0.1:8080")
	require.Error(t, err)
	require.Equal(t, code, http.StatusTooManyRequests)
}

func TestPaidTierFailFallbackToFreeTier(t *testing.T) {
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
	require.NoError(t, err)
	pub, err := info.GetPubKey()
	require.NoError(t, err)
	pk, err := common.NewPubKeyFromCrypto(pub)
	require.NoError(t, err)

	var signature []byte
	nonce := int64(3)
	service := common.BTCService.String()

	message := []byte(fmt.Sprintf("%s:%s:%s:%d", pubkey.String(), service, pk, nonce))
	signature, _, err = kb.Sign("whatever", message)
	require.NoError(t, err)

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

	contract := types.NewContract(pubkey, common.BTCService, pk)
	contract.Height = 5
	contract.Duration = 100
	contract.Id = 55556
	// set contract to expired
	proxy.MemStore.SetHeight(120)
	proxy.MemStore.Put(contract)

	nextHandler := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	})
	handlerForTest := proxy.auth(nextHandler)
	require.NotNil(t, handlerForTest)
	aa := ArkAuth{
		Nonce:      nonce,
		ContractId: contract.Id,
		Spender:    pk,
		Signature:  signature,
	}

	responseRecorder := httptest.NewRecorder()
	target := "/test?arkauth=" + aa.String()
	r := httptest.NewRequest(http.MethodPost, target, nil)
	handlerForTest.ServeHTTP(responseRecorder, r)
	require.Equal(t, responseRecorder.Code, 200)
}
