package keeper

import (
	"testing"
	"time"

	"cosmossdk.io/log"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/arkeonetwork/arkeo/x/claim/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"

	"cosmossdk.io/store"
	storemetrics "cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/arkeonetwork/arkeo/app"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

type (
	TestKeepers struct {
		ClaimKeeper   keeper.Keeper
		AccountKeeper authkeeper.AccountKeeper
		BankKeeper    bankkeeper.Keeper
	}
)

// CreateTestClaimKeepers creates test keepers for claim module
func CreateTestClaimKeepers(t testing.TB) (TestKeepers, sdk.Context) {
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	keyAcc := storetypes.NewKVStoreKey(authtypes.StoreKey)
	keyBank := storetypes.NewKVStoreKey(banktypes.StoreKey)
	keyParams := storetypes.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := storetypes.NewTransientStoreKey(paramstypes.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	logger := log.NewNopLogger()
	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, logger, storemetrics.NewNoOpMetrics())
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	encodingConfig := app.MakeEncodingConfig()
	types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cdc := utils.MakeTestMarshaler()
	amino := encodingConfig.Amino

	paramsSubspace := paramstypes.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"ClaimParams",
	)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	ctx = ctx.WithBlockTime(time.Now().UTC()) // needed for airdrop start time

	paramskeeper.NewKeeper(cdc, amino, keyParams, tkeyParams)
	govModuleAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	accountKeeper := authkeeper.NewAccountKeeper(
		cdc,
		runtime.NewKVStoreService(keyAcc),
		authtypes.ProtoBaseAccount,
		map[string][]string{
			stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
			stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
			types.ModuleName:               {authtypes.Minter},
			arkeotypes.ReserveName:         {},
			arkeotypes.ProviderName:        {},
			arkeotypes.ContractName:        {},
		},
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.Bech32PrefixAccAddr,
		govModuleAddr,
	)

	// accountKeeper.SetParams(ctx, authtypes.DefaultParams())

	bankKeeper := bankkeeper.NewBaseKeeper(
		cdc,
		runtime.NewKVStoreService(keyBank),
		accountKeeper,
		nil,
		govModuleAddr,
		logger,
	)
	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		accountKeeper,
		bankKeeper,
		memStoreKey,
		paramsSubspace,
		logger,
	)

	// Initialize params
	airdropStartTime := time.Now().UTC().Add(-time.Hour) // started an hour ago
	params := types.Params{
		AirdropStartTime:   airdropStartTime,
		DurationUntilDecay: types.DefaultDurationUntilDecay,
		DurationOfDecay:    types.DefaultDurationOfDecay,
		ClaimDenom:         types.DefaultClaimDenom,
	}

	k.SetParams(ctx, params)
	return TestKeepers{
		ClaimKeeper:   k,
		AccountKeeper: accountKeeper,
		BankKeeper:    bankKeeper,
	}, ctx
}
