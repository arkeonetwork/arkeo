package keeper

import (
	"testing"

	math "cosmossdk.io/math"
	storemetrics "cosmossdk.io/store/metrics"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/stretchr/testify/require"

	arekoappParams "github.com/arkeonetwork/arkeo/app/params"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	bech32Prefix = "tarkeo"
)

var (
	// Bech32PrefixAccAddr
	Bech32PrefixAccAddr  = bech32Prefix
	Bech32PrefixAccPub   = bech32Prefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func SetupKeeper(t testing.TB) (cosmos.Context, Keeper) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	keyAcc := cosmos.NewKVStoreKey(authtypes.StoreKey)
	keyBank := cosmos.NewKVStoreKey(banktypes.StoreKey)
	keyStake := cosmos.NewKVStoreKey(stakingtypes.StoreKey)
	keyParams := cosmos.NewKVStoreKey(typesparams.StoreKey)
	keyMint := cosmos.NewKVStoreKey(minttypes.StoreKey)
	tkeyParams := cosmos.NewTransientStoreKey(typesparams.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	keydist := cosmos.NewKVStoreKey(disttypes.StoreKey)

	cfg := sdk.GetConfig()

	cfg.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	logger := log.NewNopLogger()
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	encodingConfig := arekoappParams.MakeEncodingConfig()
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

	_ = paramskeeper.NewKeeper(cdc, amino, keyParams, tkeyParams)
	govModuleAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()
	ak := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keyAcc),
		authtypes.ProtoBaseAccount,
		map[string][]string{
			distrtypes.ModuleName:          {authtypes.Minter},
			stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
			stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
			types.ModuleName:               {authtypes.Minter, authtypes.Burner},
			types.ReserveName:              {authtypes.Minter, authtypes.Burner},
			types.ProviderName:             {},
			types.ContractName:             {},
			minttypes.ModuleName:           {authtypes.Minter},
		},
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.Bech32PrefixAccAddr,
		govModuleAddr,
	)

	bk := bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(keyBank),
		ak,
		nil,
		govModuleAddr,
		logger,
	)
	require.NoError(t, bk.SetParams(ctx, banktypes.DefaultParams()))

	sk := stakingkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keyStake),
		ak,
		bk,
		govModuleAddr,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	dk := distkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keydist),
		ak,
		bk,
		sk,
		authtypes.FeeCollectorName,
		govModuleAddr,
	)

	mk := mintkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keyMint),
		sk,
		ak,
		bk,
		authtypes.FeeCollectorName,
		govModuleAddr,
	)

	k := NewKVStore(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bk,
		ak,
		*sk,
		govModuleAddr,
		logger,
		mk,
		dk,
	)
	k.SetVersion(ctx, common.GetCurrentVersion())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return ctx, *k
}

func SetupKeeperWithStaking(t testing.TB) (cosmos.Context, Keeper, stakingkeeper.Keeper) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	keyAcc := cosmos.NewKVStoreKey(authtypes.StoreKey)
	keyBank := cosmos.NewKVStoreKey(banktypes.StoreKey)
	keyMint := cosmos.NewKVStoreKey(minttypes.StoreKey)
	keyStake := cosmos.NewKVStoreKey(stakingtypes.StoreKey)
	keyParams := cosmos.NewKVStoreKey(typesparams.StoreKey)
	tkeyParams := cosmos.NewTransientStoreKey(typesparams.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	keydist := cosmos.NewKVStoreKey(disttypes.StoreKey)

	cfg := sdk.GetConfig()

	cfg.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	cfg.SetBech32PrefixForValidator(Bech32PrefixValAddr, Bech32PrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(Bech32PrefixConsAddr, Bech32PrefixConsPub)

	logger := log.NewNopLogger()
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyStake, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(keydist, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	encodingConfig := arekoappParams.MakeEncodingConfig()
	types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cdc := utils.MakeTestMarshaler()
	amino := encodingConfig.Amino

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"ArkeoParams",
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger()).WithChainID("arkeo-1")

	govModuleAddr := "tarkeo1krj9ywwmqcellgunxg66kjw5dtt402kq0uf6pu"
	_ = paramskeeper.NewKeeper(cdc, amino, keyParams, tkeyParams)
	ak := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keyAcc),
		authtypes.ProtoBaseAccount,
		map[string][]string{
			distrtypes.ModuleName:          {authtypes.Minter},
			stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
			stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
			types.ModuleName:               {authtypes.Minter, authtypes.Burner},
			types.ReserveName:              {authtypes.Minter, authtypes.Burner},
			types.ProviderName:             {},
			types.ContractName:             {},
			govtypes.ModuleName:            {authtypes.Minter, authtypes.Burner},
			minttypes.ModuleName:           {authtypes.Minter},
		},
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.Bech32PrefixAccAddr,
		govModuleAddr,
	)

	bk := bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(keyBank),
		ak,
		nil,
		govModuleAddr,
		logger,
	)
	_ = bk.SetParams(ctx, banktypes.DefaultParams())

	sk := stakingkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keyStake),
		ak,
		bk,
		govModuleAddr,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	mk := mintkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keyMint),
		sk,
		ak,
		bk,
		authtypes.FeeCollectorName,
		govModuleAddr,
	)
	dk := distkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keydist),
		ak,
		bk,
		sk,
		authtypes.FeeCollectorName,
		govModuleAddr,
	)
	_ = dk.FeePool.Set(ctx, disttypes.FeePool{CommunityPool: []sdk.DecCoin{sdk.NewDecCoin(configs.Denom, math.NewInt(10000))}})

	k := NewKVStore(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bk,
		ak,
		*sk,
		govModuleAddr,
		logger,
		mk,
		dk,
	)
	k.SetVersion(ctx, common.GetCurrentVersion())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return ctx, *k, *sk
}
