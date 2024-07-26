package keeper

import (
	"testing"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/arkeo/keeper"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"

	// "github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func ArkeoKeeper(t testing.TB) (cosmos.Context, keeper.Keeper) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	keyAcc := cosmos.NewKVStoreKey(authtypes.StoreKey)
	keyBank := cosmos.NewKVStoreKey(banktypes.StoreKey)
	keyStake := cosmos.NewKVStoreKey(stakingtypes.StoreKey)
	keyParams := cosmos.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := cosmos.NewTransientStoreKey(paramstypes.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)
	logger := log.NewNopLogger()

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, logger, nil)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	legacyCodec := utils.MakeTestCodec()

	govModuleAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	_ = paramskeeper.NewKeeper(cdc, legacyCodec, keyParams, tkeyParams)

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"ArkeoParams",
	)
	ak := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keyAcc),
		authtypes.ProtoBaseAccount,
		map[string][]string{
			stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
			stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
			types.ModuleName:               {authtypes.Minter, authtypes.Burner},
			types.ReserveName:              {},
			types.ProviderName:             {},
			types.ContractName:             {},
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

	sk := stakingkeeper.NewKeeper(
		cdc,
		runtime.NewKVStoreService(keyStake),
		ak,
		bk,
		govModuleAddr,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)
	k := keeper.NewKVStore(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bk,
		ak,
		*sk,
	)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	k.SetVersion(ctx, common.GetCurrentVersion())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return ctx, k
}
