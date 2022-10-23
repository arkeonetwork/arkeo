package keeper

import (
	"io/ioutil"
	"mercury/common/cosmos"
	"mercury/x/mercury/keeper"
	"mercury/x/mercury/types"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/blang/semver"

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
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

// create a codec used only for testing
func MakeTestCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	banktypes.RegisterLegacyAminoCodec(cdc)
	authtypes.RegisterLegacyAminoCodec(cdc)
	types.RegisterCodec(cdc)
	cosmos.RegisterCodec(cdc)
	// codec.RegisterCrypto(cdc)
	return cdc
}

func MercuryKeeper(t testing.TB) (keeper.Keeper, sdk.Context) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := cosmos.NewKVStoreKey(authtypes.StoreKey)
	keyBank := cosmos.NewKVStoreKey(banktypes.StoreKey)
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
	legacyCodec := MakeTestCodec()

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"MercuryParams",
	)

	pk := paramskeeper.NewKeeper(cdc, legacyCodec, keyParams, tkeyParams)
	ak := authkeeper.NewAccountKeeper(cdc, keyAcc, pk.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, map[string][]string{
		types.ModuleName: {authtypes.Minter, authtypes.Burner},
	}, sdk.Bech32PrefixAccAddr)

	bk := bankkeeper.NewBaseKeeper(cdc, keyBank, ak, pk.Subspace(banktypes.ModuleName), nil)

	k := keeper.NewKVStore(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bk,
		ak,
		GetCurrentVersion(),
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return k, ctx
}

// GetCurrentVersion - intended for unit tests, fetches the current version of
// mercury via `version` file
// #nosec G304 this is a method only used for test purpose
func GetCurrentVersion() semver.Version {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "../..")
	dat, err := ioutil.ReadFile(path.Join(dir, "version"))
	if err != nil {
		panic(err)
	}
	v, err := semver.Make(strings.TrimSpace(string(dat)))
	if err != nil {
		panic(err)
	}
	return v
}
