package keeper

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/testutil/utils"
	"github.com/arkeonetwork/arkeo/x/arkeo/keeper"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func ArkeoKeeper(t testing.TB) (cosmos.Context, keeper.Keeper) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := cosmos.NewKVStoreKey(authtypes.StoreKey)
	keyBank := cosmos.NewKVStoreKey(banktypes.StoreKey)
	keyStake := cosmos.NewKVStoreKey(stakingtypes.StoreKey)
	keyParams := cosmos.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := cosmos.NewTransientStoreKey(paramstypes.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	require.NoError(t, stateStore.LoadLatestVersion())

	registry := codectypes.NewInterfaceRegistry()
	cdc := codec.NewProtoCodec(registry)
	legacyCodec := utils.MakeTestCodec()

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"ArkeoParams",
	)

	pk := paramskeeper.NewKeeper(cdc, legacyCodec, keyParams, tkeyParams)
	ak := authkeeper.NewAccountKeeper(cdc, keyAcc, pk.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, map[string][]string{
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		types.ModuleName:               {authtypes.Minter, authtypes.Burner},
		types.ReserveName:              {},
		types.ProviderName:             {},
		types.ContractName:             {},
	}, sdk.Bech32PrefixAccAddr)

	bk := bankkeeper.NewBaseKeeper(cdc, keyBank, ak, pk.Subspace(banktypes.ModuleName), nil)

	sk := stakingkeeper.NewKeeper(cdc, keyStake, ak, bk, pk.Subspace(stakingtypes.ModuleName))
	k := keeper.NewKVStore(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bk,
		ak,
		sk,
	)
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	k.SetStoreVersion(ctx, common.GetCurrentVersion())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return ctx, k
}
