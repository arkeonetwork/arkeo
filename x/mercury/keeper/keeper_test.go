package keeper

import (
	"errors"
	"fmt"
	"mercury/common/cosmos"
	"mercury/x/mercury/types"
	"testing"

	"github.com/blang/semver"
	. "gopkg.in/check.v1"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	simappparams "github.com/cosmos/cosmos-sdk/simapp/params"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/cosmos/cosmos-sdk/x/capability"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/evidence"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/mint"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	typesparams "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/upgrade"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func Test(t *testing.T) { TestingT(t) }

// ModuleBasics is a mock module basic manager for testing
var ModuleBasics = module.NewBasicManager(
	auth.AppModuleBasic{},
	genutil.AppModuleBasic{},
	bank.AppModuleBasic{},
	capability.AppModuleBasic{},
	staking.AppModuleBasic{},
	mint.AppModuleBasic{},
	distribution.AppModuleBasic{},
	/*
		gov.NewAppModuleBasic(
			paramsclient.ProposalHandler, distrclient.ProposalHandler, upgradeclient.ProposalHandler, upgradeclient.CancelProposalHandler,
		),
	*/
	params.AppModuleBasic{},
	crisis.AppModuleBasic{},
	slashing.AppModuleBasic{},
	upgrade.AppModuleBasic{},
	evidence.AppModuleBasic{},
	vesting.AppModuleBasic{},
)

func MakeTestMarshaler() codec.Codec {
	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	ModuleBasics.RegisterInterfaces(interfaceRegistry)
	types.RegisterInterfaces(interfaceRegistry)
	return codec.NewProtoCodec(interfaceRegistry)
}

func SetupKeeper(c *C) (cosmos.Context, Keeper) {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	keyAcc := cosmos.NewKVStoreKey(authtypes.StoreKey)
	keyBank := cosmos.NewKVStoreKey(banktypes.StoreKey)
	keyParams := cosmos.NewKVStoreKey(paramstypes.StoreKey)
	tkeyParams := cosmos.NewTransientStoreKey(paramstypes.TStoreKey)
	memStoreKey := storetypes.NewMemoryStoreKey(types.MemStoreKey)

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)
	stateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyAcc, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyBank, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(keyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(tkeyParams, storetypes.StoreTypeIAVL, db)
	stateStore.MountStoreWithDB(memStoreKey, storetypes.StoreTypeMemory, nil)
	c.Assert(stateStore.LoadLatestVersion(), IsNil)

	encodingConfig := simappparams.MakeTestEncodingConfig()
	types.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	cdc := MakeTestMarshaler()
	amino := encodingConfig.Amino

	paramsSubspace := typesparams.NewSubspace(cdc,
		types.Amino,
		storeKey,
		memStoreKey,
		"MercuryParams",
	)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

	pk := paramskeeper.NewKeeper(cdc, amino, keyParams, tkeyParams)
	ak := authkeeper.NewAccountKeeper(cdc, keyAcc, pk.Subspace(authtypes.ModuleName), authtypes.ProtoBaseAccount, map[string][]string{
		types.ModuleName:  {authtypes.Minter, authtypes.Burner},
		types.ReserveName: {},
	}, sdk.Bech32PrefixAccAddr)
	ak.SetParams(ctx, authtypes.DefaultParams())

	bk := bankkeeper.NewBaseKeeper(cdc, keyBank, ak, pk.Subspace(banktypes.ModuleName), nil)
	bk.SetParams(ctx, banktypes.DefaultParams())

	k := NewKVStore(
		cdc,
		storeKey,
		memStoreKey,
		paramsSubspace,
		bk,
		ak,
		semver.MustParse("0.0.0"),
	)

	// Initialize params
	k.SetParams(ctx, types.DefaultParams())

	return ctx, *k
}

type errIsChecker struct {
	*CheckerInfo
}

var ErrIs Checker = &errIsChecker{
	&CheckerInfo{Name: "ErrIs", Params: []string{"obtained", "expected"}},
}

func (errIsChecker) Check(params []interface{}, names []string) (result bool, err string) {
	result = errors.Is(params[0].(error), params[1].(error))
	if !result {
		err = fmt.Sprintf("Errors do not match!\nObtained: %s\nExpected: %s", params[0].(error), params[1].(error))
	}
	return
}
