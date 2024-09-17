package keeper

import (
	"context"
	"fmt"
	"strings"

	"cosmossdk.io/errors"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type dbPrefix string

func (p dbPrefix) String() string {
	return string(p)
}

type Keeper interface {
	Logger(ctx sdk.Context) log.Logger
	GetParams(ctx sdk.Context) types.Params
	SetParams(ctx sdk.Context, params types.Params)
	Cdc() codec.BinaryCodec
	GetComputedVersion(ctx cosmos.Context) int64
	GetVersion(ctx cosmos.Context) int64
	SetVersion(ctx cosmos.Context, ver int64)
	GetKey(ctx cosmos.Context, prefix dbPrefix, key string) string
	GetVersionForAddress(ctx cosmos.Context, _ cosmos.ValAddress) int64
	SetVersionForAddress(ctx cosmos.Context, _ cosmos.ValAddress, ver int64)
	GetSupply(ctx cosmos.Context, denom string) cosmos.Coin
	GetBalanceOfModule(ctx cosmos.Context, moduleName, denom string) cosmos.Int
	SendFromModuleToModule(ctx cosmos.Context, from, to string, coin cosmos.Coins) error
	SendFromAccountToModule(ctx cosmos.Context, from cosmos.AccAddress, to string, _ cosmos.Coins) error
	SendFromModuleToAccount(ctx cosmos.Context, from string, to cosmos.AccAddress, _ cosmos.Coins) error
	MintToModule(ctx cosmos.Context, module string, coin cosmos.Coin) error
	BurnFromModule(ctx cosmos.Context, module string, coin cosmos.Coin) error
	MintAndSendToAccount(ctx cosmos.Context, to cosmos.AccAddress, coin cosmos.Coin) error
	GetModuleAccAddress(module string) cosmos.AccAddress
	GetBalance(ctx cosmos.Context, addr cosmos.AccAddress) cosmos.Coins
	HasCoins(ctx cosmos.Context, addr cosmos.AccAddress, coins cosmos.Coins) bool

	// passthrough funcs
	SendCoins(ctx cosmos.Context, from, to cosmos.AccAddress, coins cosmos.Coins) error
	AddCoins(ctx cosmos.Context, addr cosmos.AccAddress, coins cosmos.Coins) error
	GetActiveValidators(ctx cosmos.Context) ([]stakingtypes.Validator, error)
	GetAccount(ctx cosmos.Context, addr cosmos.AccAddress) cosmos.Account
	StakingSetParams(ctx cosmos.Context, params stakingtypes.Params) error
	MintAndDistributeTokens(ctx cosmos.Context, newlyMinted cosmos.Coin) (cosmos.Coin, error)
	GetCirculatingSupply(ctx cosmos.Context, denom string) (cosmos.Coin, error)
	GetInflationRate(ctx cosmos.Context) (math.LegacyDec, error)

	// Query
	Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error)
	FetchProvider(c context.Context, req *types.QueryFetchProviderRequest) (*types.QueryFetchProviderResponse, error)
	ProviderAll(c context.Context, req *types.QueryAllProviderRequest) (*types.QueryAllProviderResponse, error)
	FetchContract(c context.Context, req *types.QueryFetchContractRequest) (*types.QueryFetchContractResponse, error)
	ContractAll(c context.Context, req *types.QueryAllContractRequest) (*types.QueryAllContractResponse, error)
	ActiveContract(goCtx context.Context, req *types.QueryActiveContractRequest) (*types.QueryActiveContractResponse, error)

	// Keeper Interfaces
	KeeperProvider
	KeeperContract
}

type KeeperProvider interface {
	GetProviderIterator(_ cosmos.Context) cosmos.Iterator
	GetProvider(_ cosmos.Context, _ common.PubKey, _ common.Service) (types.Provider, error)
	SetProvider(_ cosmos.Context, _ types.Provider) error
	ProviderExists(_ cosmos.Context, _ common.PubKey, _ common.Service) bool
	RemoveProvider(_ cosmos.Context, _ common.PubKey, _ common.Service)
}

type KeeperContract interface {
	GetContractIterator(_ cosmos.Context) cosmos.Iterator
	GetContract(_ cosmos.Context, _ uint64) (types.Contract, error)
	SetContract(_ cosmos.Context, _ types.Contract) error
	ContractExists(_ cosmos.Context, _ uint64) bool
	RemoveContract(_ cosmos.Context, _ uint64)
	GetContractExpirationSetIterator(_ cosmos.Context) cosmos.Iterator
	GetUserContractSetIterator(_ cosmos.Context) cosmos.Iterator
	GetContractExpirationSet(_ cosmos.Context, _ int64) (types.ContractExpirationSet, error)
	SetContractExpirationSet(_ cosmos.Context, _ types.ContractExpirationSet) error
	RemoveContractExpirationSet(_ cosmos.Context, _ int64)
	RemoveFromUserContractSet(ctx cosmos.Context, user common.PubKey, contractId uint64) error
	GetNextContractId(_ cosmos.Context) uint64
	SetNextContractId(ctx cosmos.Context, contractId uint64)
	GetAndIncrementNextContractId(ctx cosmos.Context) uint64
	SetUserContractSet(ctx cosmos.Context, contractSet types.UserContractSet) error
	GetUserContractSet(ctx cosmos.Context, pubkey common.PubKey) (types.UserContractSet, error)
	GetActiveContractForUser(ctx cosmos.Context, user, provider common.PubKey, service common.Service) (types.Contract, error)
}

const (
	prefixVersion               dbPrefix = "ver/"
	prefixProvider              dbPrefix = "p/"
	prefixContract              dbPrefix = "c/"
	prefixContractNextId        dbPrefix = "cni/"
	prefixContractExpirationSet dbPrefix = "ces/"
	prefixUserContractSet       dbPrefix = "ucs/"
)

type KVStore struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	memKey        storetypes.StoreKey
	paramstore    paramtypes.Subspace
	coinKeeper    bankkeeper.Keeper
	accountKeeper authkeeper.AccountKeeper
	stakingKeeper stakingkeeper.Keeper
	authority     string
	logger        log.Logger
	mintKeeper    minttypes.Keeper
}

func NewKVStore(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,

	/*
			accountKeeper       keeper.AccountKeeper
		    bankKeeper          keeper.Keeper
		    distributionKeeper  distkeeper.Keeper
		    stakingKeeper       stakingkeeper.Keeper
	*/
	coinKeeper bankkeeper.Keeper,
	accountKeeper authkeeper.AccountKeeper,
	stakingKeeper stakingkeeper.Keeper,
	authority string,
	logger log.Logger,
	mintKeeper minttypes.Keeper,
) *KVStore {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &KVStore{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramstore:    ps,
		coinKeeper:    coinKeeper,
		accountKeeper: accountKeeper,
		stakingKeeper: stakingKeeper,
		authority:     authority,
		logger:        logger,
		mintKeeper:    mintKeeper,
	}
}

func (k KVStore) Logger(ctx sdk.Context) log.Logger {
	return k.logger.With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetKey return a key that can be used to store into key value store
func (k KVStore) GetKey(ctx cosmos.Context, prefix dbPrefix, key string) string {
	return fmt.Sprintf("%s/%s", prefix, strings.ToUpper(key))
}

// Cdc return the amino codec
func (k KVStore) Cdc() codec.BinaryCodec {
	return k.cdc
}

// GetParams get all parameters as types.Params
func (k KVStore) GetParams(ctx sdk.Context) types.Params {
	return types.NewParams()
}

// SetParams set the params
func (k KVStore) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}

// TODO: Check Thi Again
func (k KVStore) GetComputedVersion(ctx cosmos.Context) int64 {
	versions := make(map[int64]int64) // maps are safe in blockchains, but should be okay in this case
	validators, err := k.stakingKeeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		k.Logger(ctx).Error(err.Error())
	}

	// if there is only one validator, no need for consensus. Just return the
	// validator's current version. This also helps makes
	// integration/regression tests run the latest version
	if len(validators) == 1 {
		return configs.SWVersion
	}

	currentVersion := k.GetVersion(ctx)
	minNum := configs.GetConfigValues(currentVersion).GetInt64Value(configs.VersionConsensus)
	min := int64(len(validators)) * minNum / 100

	for _, val := range validators {
		if !val.IsBonded() {
			continue
		}

		valBz, err := k.stakingKeeper.ValidatorAddressCodec().StringToBytes(val.GetOperator())
		if err != nil {
			k.Logger(ctx).Error(err.Error())
		}
		ver := k.GetVersionForAddress(ctx, valBz)
		if _, ok := versions[ver]; !ok {
			versions[ver] = 0
		}
		versions[ver] += 1
		if versions[ver] >= min {
			return ver
		}
	}
	return currentVersion
}

// SetVersion save the store version
func (k KVStore) SetVersion(ctx cosmos.Context, value int64) {
	key := k.GetKey(ctx, prefixVersion, "")
	store := ctx.KVStore(k.storeKey)
	ver := types.ProtoInt64{Value: value}
	store.Set([]byte(key), k.cdc.MustMarshal(&ver))
}

// GetVersion get the current key value store version
func (k KVStore) GetVersion(ctx cosmos.Context) int64 {
	key := k.GetKey(ctx, prefixVersion, "")
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return 1
	}
	var ver types.ProtoInt64
	buf := store.Get([]byte(key))
	k.cdc.MustUnmarshal(buf, &ver)
	return ver.Value
}

// SetVersionForAddress save the store version
func (k KVStore) SetVersionForAddress(ctx cosmos.Context, addr cosmos.ValAddress, value int64) {
	key := k.GetKey(ctx, prefixVersion, addr.String())
	store := ctx.KVStore(k.storeKey)
	ver := types.ProtoInt64{Value: value}
	store.Set([]byte(key), k.cdc.MustMarshal(&ver))
}

// GetVersionForAddress get the current key value store version
func (k KVStore) GetVersionForAddress(ctx cosmos.Context, addr cosmos.ValAddress) int64 {
	key := k.GetKey(ctx, prefixVersion, addr.String())
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(key)) {
		return 1
	}
	var ver types.ProtoInt64
	buf := store.Get([]byte(key))
	k.cdc.MustUnmarshal(buf, &ver)
	return ver.Value
}

// getIterator - get an iterator for given prefix
func (k KVStore) getIterator(ctx cosmos.Context, prefix dbPrefix) cosmos.Iterator {
	store := ctx.KVStore(k.storeKey)
	return cosmos.KVStorePrefixIterator(store, []byte(prefix))
}

// del - delete data from the kvstore
func (k KVStore) del(ctx cosmos.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(key)) {
		store.Delete([]byte(key))
	}
}

// has - kvstore has key
func (k KVStore) has(ctx cosmos.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(key))
}

func (k KVStore) GetBalanceOfModule(ctx cosmos.Context, moduleName, denom string) cosmos.Int {
	addr := k.accountKeeper.GetModuleAddress(moduleName)
	coin := k.coinKeeper.GetBalance(ctx, addr, denom)
	return cosmos.NewIntFromBigInt(coin.Amount.BigInt())
}

func (k KVStore) GetSupply(ctx cosmos.Context, denom string) cosmos.Coin {
	return k.coinKeeper.GetSupply(ctx, denom)
}

// SendFromModuleToModule transfer asset from one module to another
func (k KVStore) SendFromModuleToModule(ctx cosmos.Context, from, to string, coins cosmos.Coins) error {
	return k.coinKeeper.SendCoinsFromModuleToModule(ctx, from, to, coins)
}

func (k KVStore) SendCoins(ctx cosmos.Context, from, to cosmos.AccAddress, coins cosmos.Coins) error {
	return k.coinKeeper.SendCoins(ctx, from, to, coins)
}

func (k KVStore) AddCoins(ctx cosmos.Context, addr cosmos.AccAddress, coins cosmos.Coins) error {
	return k.coinKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr, coins)
}

// SendFromAccountToModule transfer fund from one account to a module
func (k KVStore) SendFromAccountToModule(ctx cosmos.Context, from cosmos.AccAddress, to string, coins cosmos.Coins) error {
	if !k.HasCoins(ctx, from, coins) {
		return errors.Wrapf(sdkerrors.ErrInsufficientFunds, "not enough balance for account %s", from)
	}
	return k.coinKeeper.SendCoinsFromAccountToModule(ctx, from, to, coins)
}

// SendFromModuleToAccount transfer fund from module to an account
func (k KVStore) SendFromModuleToAccount(ctx cosmos.Context, from string, to cosmos.AccAddress, coins cosmos.Coins) error {
	return k.coinKeeper.SendCoinsFromModuleToAccount(ctx, from, to, coins)
}

func (k KVStore) BurnFromModule(ctx cosmos.Context, module string, coin cosmos.Coin) error {
	return k.coinKeeper.BurnCoins(ctx, module, cosmos.Coins{coin})
}

func (k KVStore) MintToModule(ctx cosmos.Context, module string, coin cosmos.Coin) error {
	return k.coinKeeper.MintCoins(ctx, module, cosmos.Coins{coin})
}

func (k KVStore) MintAndSendToAccount(ctx cosmos.Context, to cosmos.AccAddress, coin cosmos.Coin) error {
	// Mint coins into the reserve
	if err := k.MintToModule(ctx, types.ModuleName, coin); err != nil {
		return err
	}
	return k.SendFromModuleToAccount(ctx, types.ModuleName, to, cosmos.NewCoins(coin))
}

func (k KVStore) GetModuleAddress(module string) (cosmos.AccAddress, error) {
	return cosmos.AccAddressFromBech32(module)
}

func (k KVStore) GetModuleAccAddress(module string) cosmos.AccAddress {
	return k.accountKeeper.GetModuleAddress(module)
}

func (k KVStore) GetBalance(ctx cosmos.Context, addr cosmos.AccAddress) cosmos.Coins {
	return k.coinKeeper.GetAllBalances(ctx, addr)
}

func (k KVStore) HasCoins(ctx cosmos.Context, addr cosmos.AccAddress, coins cosmos.Coins) bool {
	balance := k.coinKeeper.GetAllBalances(ctx, addr)
	return balance.IsAllGTE(coins)
}

func (k KVStore) GetAccount(ctx cosmos.Context, addr cosmos.AccAddress) cosmos.Account {
	return k.accountKeeper.GetAccount(ctx, addr)
}

func (k KVStore) GetActiveValidators(ctx cosmos.Context) ([]stakingtypes.Validator, error) {
	return k.stakingKeeper.GetBondedValidatorsByPower(ctx)
}

func (k KVStore) StakingSetParams(ctx cosmos.Context, params stakingtypes.Params) error {
	return k.stakingKeeper.SetParams(ctx, params)
}

func (k KVStore) GetAuthority() string {
	return k.authority
}

func (k KVStore) GetCirculatingSupply(ctx cosmos.Context, denom string) (cosmos.Coin, error) {

	sdkContext := sdk.UnwrapSDKContext(ctx)

	// Get Total Supply
	fullTokenSupply, err := k.coinKeeper.TotalSupply(ctx, &banktypes.QueryTotalSupplyRequest{})
	if err != nil {
		sdkContext.Logger().Error("Failed to get full token supply data", err)
		return cosmos.NewCoin(denom, sdkmath.NewInt(0)), err
	}
	totalSupply := fullTokenSupply.Supply.AmountOf(configs.Denom)
	sdkContext.Logger().Error("Current supply of token with denom :", denom, "supply:", totalSupply)

	// Module Address for which the circulating supply should be exempted
	modulesToExempt := []string{types.GrantPool, types.CommunityPool, types.DevFundPool}
	exemptBalance := cosmos.NewInt(0)

	// range over the module and create exempt balances
	for _, moduleName := range modulesToExempt {
		moduleAddr := k.accountKeeper.GetModuleAddress(moduleName)
		moduleBalance := k.coinKeeper.GetBalance(ctx, moduleAddr, denom)
		exemptBalance = exemptBalance.Add(moduleBalance.Amount)
	}

	// total supply without balances of exempted module
	circulatingSupply := totalSupply.Sub(exemptBalance)

	return cosmos.NewCoin(denom, circulatingSupply), nil
}

func (k KVStore) MintAndDistributeTokens(ctx cosmos.Context, newlyMinted cosmos.Coin) (cosmos.Coin, error) {

	params := k.GetParams(ctx)

	newlyMintedDec := sdkmath.LegacyNewDecFromInt(newlyMinted.Amount)

	devFundAmount := newlyMintedDec.Mul(params.DevFundPercentage).TruncateInt()
	communityPoolAmount := newlyMintedDec.Mul(params.CommunityPoolPercentage).TruncateInt()
	// validatorRewardAmount := newlyMintedDec.Mul(params.ValidatorRewardsPercentage).TruncateInt()

	if err := k.MintToModule(ctx, types.DevFundPool, cosmos.NewCoin(newlyMinted.Denom, devFundAmount)); err != nil {
		return cosmos.NewCoin(configs.Denom, sdkmath.NewInt(0)), fmt.Errorf("error sending amount to module %s", err)
	}

	if err := k.MintToModule(ctx, types.CommunityPool, cosmos.NewCoin(newlyMinted.Denom, communityPoolAmount)); err != nil {
		return cosmos.NewCoin(configs.Denom, sdkmath.NewInt(0)), fmt.Errorf("error sending amount to module %s", err)
	}
	balance := newlyMintedDec.Sub(sdkmath.LegacyDec(devFundAmount)).Sub(sdkmath.LegacyDec(communityPoolAmount))

	return cosmos.NewCoin(configs.Denom, balance.RoundInt()), nil

}

func (k KVStore) GetInflationRate(ctx cosmos.Context) (math.LegacyDec, error) {
	minter, err := k.mintKeeper.Minter.Get(ctx)
	if err != nil {
		return math.LegacyNewDec(0), err
	}

	return minter.Inflation, nil
}
