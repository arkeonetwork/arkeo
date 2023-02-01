package keeper

import (
	"context"
	"errors"
	"fmt"

	"github.com/ArkeoNetwork/arkeo-protocol/common"
	"github.com/ArkeoNetwork/arkeo-protocol/common/cosmos"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"

	"github.com/blang/semver"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"
)

var kaboom = errors.New("Kaboom!!!")

type KVStoreDummy struct{}

func (k KVStoreDummy) Cdc() codec.BinaryCodec {
	return nil
}

func (k KVStoreDummy) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams()
}
func (k KVStoreDummy) SetParams(ctx sdk.Context, params types.Params) {}
func (k KVStoreDummy) CoinKeeper() bankkeeper.Keeper                  { return bankkeeper.BaseKeeper{} }
func (k KVStoreDummy) AccountKeeper() authkeeper.AccountKeeper        { return authkeeper.AccountKeeper{} }
func (k KVStoreDummy) Logger(ctx cosmos.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k KVStoreDummy) GetVersion() semver.Version { return semver.MustParse("9999999.0.0") }
func (k KVStoreDummy) GetVersionWithCtx(ctx cosmos.Context) (semver.Version, bool) {
	return semver.MustParse("9999999.0.0"), true
}
func (k KVStoreDummy) SetVersionWithCtx(ctx cosmos.Context, v semver.Version) {}

func (k KVStoreDummy) GetKey(_ cosmos.Context, prefix dbPrefix, key string) string {
	return fmt.Sprintf("%s/1/%s", prefix, key)
}

func (k KVStoreDummy) GetStoreVersion(ctx cosmos.Context) int64      { return 1 }
func (k KVStoreDummy) SetStoreVersion(ctx cosmos.Context, ver int64) {}

func (k KVStoreDummy) GetRuneBalanceOfModule(ctx cosmos.Context, moduleName string) cosmos.Int {
	return cosmos.ZeroInt()
}

func (k KVStoreDummy) GetBalanceOfModule(ctx cosmos.Context, moduleName, denom string) cosmos.Int {
	return cosmos.ZeroInt()
}

func (k KVStoreDummy) SendFromModuleToModule(ctx cosmos.Context, from, to string, coins cosmos.Coins) error {
	return kaboom
}

func (k KVStoreDummy) SendCoins(ctx cosmos.Context, from, to cosmos.AccAddress, coins cosmos.Coins) error {
	return kaboom
}

func (k KVStoreDummy) AddCoins(ctx cosmos.Context, _ cosmos.AccAddress, coins cosmos.Coins) error {
	return kaboom
}

func (k KVStoreDummy) SendFromAccountToModule(ctx cosmos.Context, from cosmos.AccAddress, to string, coins cosmos.Coins) error {
	return kaboom
}

func (k KVStoreDummy) SendFromModuleToAccount(ctx cosmos.Context, from string, to cosmos.AccAddress, coins cosmos.Coins) error {
	return kaboom
}

func (k KVStoreDummy) MintToModule(ctx cosmos.Context, module string, coin cosmos.Coin) error {
	return kaboom
}

func (k KVStoreDummy) BurnFromModule(ctx cosmos.Context, module string, coin cosmos.Coin) error {
	return kaboom
}

func (k KVStoreDummy) MintAndSendToAccount(ctx cosmos.Context, to cosmos.AccAddress, coin cosmos.Coin) error {
	return kaboom
}

func (k KVStoreDummy) GetModuleAddress(module string) cosmos.AccAddress {
	return cosmos.AccAddress{}
}

func (k KVStoreDummy) GetModuleAccAddress(module string) cosmos.AccAddress {
	return nil
}

func (k KVStoreDummy) GetAccount(ctx cosmos.Context, addr cosmos.AccAddress) cosmos.Account {
	return nil
}

func (k KVStoreDummy) GetBalance(ctx cosmos.Context, addr cosmos.AccAddress) cosmos.Coins {
	return nil
}

func (k KVStoreDummy) HasCoins(ctx cosmos.Context, addr cosmos.AccAddress, coins cosmos.Coins) bool {
	return false
}

func (k KVStoreDummy) GetActiveValidators(ctx cosmos.Context) []stakingtypes.Validator {
	return nil
}

func (k KVStoreDummy) GetProviderIterator(_ cosmos.Context) cosmos.Iterator { return nil }
func (k KVStoreDummy) GetProvider(_ cosmos.Context, _ common.PubKey, _ common.Chain) (types.Provider, error) {
	return types.Provider{}, kaboom
}
func (k KVStoreDummy) SetProvider(_ cosmos.Context, _ types.Provider) error { return kaboom }
func (k KVStoreDummy) ProviderExists(_ cosmos.Context, _ common.PubKey, _ common.Chain) bool {
	return false
}
func (k KVStoreDummy) RemoveProvider(_ cosmos.Context, _ common.PubKey, _ common.Chain) {}
func (k KVStoreDummy) GetContractIterator(_ cosmos.Context) cosmos.Iterator             { return nil }
func (k KVStoreDummy) GetContract(_ cosmos.Context, _ common.PubKey, _ common.Chain, _ common.PubKey) (types.Contract, error) {
	return types.Contract{}, kaboom
}
func (k KVStoreDummy) SetContract(_ cosmos.Context, _ types.Contract) error { return kaboom }
func (k KVStoreDummy) ContractExists(_ cosmos.Context, _ common.PubKey, _ common.Chain, _ common.PubKey) bool {
	return false
}

func (k KVStoreDummy) RemoveContract(_ cosmos.Context, _ common.PubKey, _ common.Chain, _ common.PubKey) {
}

func (k KVStoreDummy) GetContractExpirationSetIterator(_ cosmos.Context) cosmos.Iterator { return nil }
func (k KVStoreDummy) GetContractExpirationSet(_ cosmos.Context, _ int64) (types.ContractExpirationSet, error) {
	return types.ContractExpirationSet{}, kaboom
}

func (k KVStoreDummy) SetContractExpirationSet(_ cosmos.Context, _ types.ContractExpirationSet) error {
	return kaboom
}
func (k KVStoreDummy) RemoveContractExpirationSet(_ cosmos.Context, _ int64) {}

func (k KVStoreDummy) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	return nil, kaboom
}

func (k KVStoreDummy) FetchProvider(c context.Context, req *types.QueryFetchProviderRequest) (*types.QueryFetchProviderResponse, error) {
	return nil, kaboom
}

func (k KVStoreDummy) ProviderAll(c context.Context, req *types.QueryAllProviderRequest) (*types.QueryAllProviderResponse, error) {
	return nil, kaboom
}

func (k KVStoreDummy) FetchContract(c context.Context, req *types.QueryFetchContractRequest) (*types.QueryFetchContractResponse, error) {
	return nil, kaboom
}

func (k KVStoreDummy) ContractAll(c context.Context, req *types.QueryAllContractRequest) (*types.QueryAllContractResponse, error) {
	return nil, kaboom
}

func (k KVStoreDummy) StakingSetParams(ctx cosmos.Context, params stakingtypes.Params) {}
