package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"

	"github.com/blang/semver"

	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func SetupKeeper(t testing.TB) (cosmos.Context, Keeper) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := cosmos.NewKVStoreKey(authtypes.StoreKey)
	keyBank := cosmos.NewKVStoreKey(banktypes.StoreKey)
	keyStake := cosmos.NewKVStoreKey(stakingtypes.StoreKey)
	keyParams := cosmos.NewKVStoreKey(typesparams.StoreKey)
	tkeyParams := cosmos.NewTransientStoreKey(typesparams.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	encodingConfig := simappparams.MakeTestEncodingConfig()
	types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cdc := utils.MakeTestMarshaler()
	amino := encodingConfig.Amino

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"ArkeoParams",
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	pk := paramskeeper.NewKeeper(cdc, amino, keyParams, tkeyParams)
	ak := authkeeper.NewAccountKeeper(cdc, keyAcc, pk.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, map[string][]string{
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
		types.ReserveName:              {},
		types.ProviderName:             {},
		types.ContractName:             {},
	}, sdk.Bech32PrefixAccAddr)
	ak.SetParams(ctx, authtypes.DefaultParams())

	bk := bankkeeper.NewBaseKeeper(cdc, keyBank, ak, pk.Subspace(banktypes.ModuleName), nil)
	bk.SetParams(ctx, banktypes.DefaultParams())

	sk := stakingkeeper.NewKeeper(cdc, keyStake, ak, bk, pk.Subspace(stakingtypes.ModuleName))
	k := NewKVStore(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bk,
		ak,
		sk,
		semver.MustParse("0.0.0"),
	)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return ctx, *k
}

func SetupKeeperWithStaking(t testing.TB) (cosmos.Context, Keeper, stakingkeeper.Keeper) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := cosmos.NewKVStoreKey(authtypes.StoreKey)
	keyBank := cosmos.NewKVStoreKey(banktypes.StoreKey)
	keyStake := cosmos.NewKVStoreKey(stakingtypes.StoreKey)
	keyParams := cosmos.NewKVStoreKey(typesparams.StoreKey)
	tkeyParams := cosmos.NewTransientStoreKey(typesparams.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyStake, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	encodingConfig := simappparams.MakeTestEncodingConfig()
	types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cdc := utils.MakeTestMarshaler()
	amino := encodingConfig.Amino

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"ArkeoParams",
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	pk := paramskeeper.NewKeeper(cdc, amino, keyParams, tkeyParams)
	ak := authkeeper.NewAccountKeeper(cdc, keyAcc, pk.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, map[string][]string{
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
		types.ReserveName:              {},
		types.ProviderName:             {},
		types.ContractName:             {},
	}, sdk.Bech32PrefixAccAddr)
	ak.SetParams(ctx, authtypes.DefaultParams())

	bk := bankkeeper.NewBaseKeeper(cdc, keyBank, ak, pk.Subspace(banktypes.ModuleName), nil)
	bk.SetParams(ctx, banktypes.DefaultParams())

	sk := stakingkeeper.NewKeeper(cdc, keyStake, ak, bk, pk.Subspace(stakingtypes.ModuleName))
	sk.SetParams(ctx, stakingtypes.DefaultParams())

	k := NewKVStore(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bk,
		ak,
		sk,
		semver.MustParse("0.0.0"),
	)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return ctx, *k, sk
}
