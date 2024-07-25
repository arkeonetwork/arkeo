package keeper

import (
	"testing"
	"time"

	"github.com/arkeonetwork/arkeo/testutil/utils"
	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/arkeonetwork/arkeo/x/claim/keeper"
	"github.com/arkeonetwork/arkeo/x/claim/types"

	"cosmossdk.io/store"
	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmdb "github.com/cosmos/cosmos-db"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
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
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := sdk.NewKVStoreKey(authtypes.StoreKey)
	keyBank := sdk.NewKVStoreKey(banktypes.StoreKey)
	keyParams := sdk.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	stateStore.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	require.NoError(t, stateStore.LoadLatestVersion())

	encodingConfig := MakeTestEncodingConfig()
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

	paramsKeeper := paramskeeper.NewKeeper(cdc, amino, keyParams, tkeyParams)
	accountKeeper := authkeeper.NewAccountKeeper(cdc, keyAcc, paramsKeeper.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, map[string][]string{
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		types.ModuleName:               {authtypes.Minter},
		arkeotypes.ReserveName:         {},
		arkeotypes.ProviderName:        {},
		arkeotypes.ContractName:        {},
	}, sdk.Bech32PrefixAccAddr)
	accountKeeper.SetParams(ctx, authtypes.DefaultParams())
	bankKeeper := bankkeeper.NewBaseKeeper(cdc, keyBank, accountKeeper, paramsKeeper.Subspace(banktypes.ModuleName), nil)
	bankKeeper.SetParams(ctx, banktypes.DefaultParams())

	k := keeper.NewKeeper(
		cdc,
		storeKey,
		accountKeeper,
		bankKeeper,
		memStoreKey,
		paramsSubspace,
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
