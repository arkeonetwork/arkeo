//go:build !regtest
// +build !regtest

package app

import (
	context "context"
	"fmt"
	"io"
	"os"

	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	reflectionv1 "cosmossdk.io/api/cosmos/reflection/v1"
	clienthelpers "cosmossdk.io/client/v2/helpers"
	"cosmossdk.io/math"
	runtimeservices "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/gogoproto/proto"

	"cosmossdk.io/client/v2/autocli"
	"cosmossdk.io/core/appmodule"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	nodeservice "github.com/cosmos/cosmos-sdk/client/grpc/node"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	runtimeservice "github.com/cosmos/cosmos-sdk/runtime/services"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	authcodec "github.com/cosmos/cosmos-sdk/x/auth/codec"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	"github.com/cosmos/cosmos-sdk/x/auth/posthandler"
	authsims "github.com/cosmos/cosmos-sdk/x/auth/simulation"
	"github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtx "github.com/cosmos/cosmos-sdk/x/auth/tx"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/cosmos-sdk/x/auth/vesting"
	vestingtypes "github.com/cosmos/cosmos-sdk/x/auth/vesting/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	authzmodule "github.com/cosmos/cosmos-sdk/x/authz/module"
	"github.com/cosmos/cosmos-sdk/x/bank"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	_ "github.com/cosmos/cosmos-sdk/x/consensus" // import for side-effects
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"
	"github.com/cosmos/cosmos-sdk/x/crisis"
	crisiskeeper "github.com/cosmos/cosmos-sdk/x/crisis/keeper"
	crisistypes "github.com/cosmos/cosmos-sdk/x/crisis/types"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"

	"github.com/arkeonetwork/arkeo/app/keepers"
	arekoappParams "github.com/arkeonetwork/arkeo/app/params"

	// distrclient "github.com/cosmos/cosmos-sdk/x/distribution/client"
	"cosmossdk.io/x/evidence"
	evidencekeeper "cosmossdk.io/x/evidence/keeper"
	evidencetypes "cosmossdk.io/x/evidence/types"
	"cosmossdk.io/x/feegrant"
	feegrantkeeper "cosmossdk.io/x/feegrant/keeper"
	feegrantmodule "cosmossdk.io/x/feegrant/module"
	"cosmossdk.io/x/upgrade"
	consensusparamtypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	genutiltypes "github.com/cosmos/cosmos-sdk/x/genutil/types"
	"github.com/cosmos/cosmos-sdk/x/gov"
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
	govkeeper "github.com/cosmos/cosmos-sdk/x/gov/keeper"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"github.com/cosmos/cosmos-sdk/x/group"
	groupkeeper "github.com/cosmos/cosmos-sdk/x/group/keeper"
	groupmodule "github.com/cosmos/cosmos-sdk/x/group/module"
	"github.com/cosmos/cosmos-sdk/x/mint"
	mintkeeper "github.com/cosmos/cosmos-sdk/x/mint/keeper"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	paramsclient "github.com/cosmos/cosmos-sdk/x/params/client"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	paramproposal "github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	slashingkeeper "github.com/cosmos/cosmos-sdk/x/slashing/keeper"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// upgradeclient "cosmossdk.io/x/upgrade/client"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	ica "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts"
	icahost "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host"

	icatypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/types"
	"github.com/cosmos/ibc-go/v8/modules/apps/transfer"

	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	// ibcclientclient "github.com/cosmos/ibc-go/v8/modules/core/02-client/client"

	"cosmossdk.io/log"
	_ "cosmossdk.io/x/circuit" // import for side-effects
	abci "github.com/cometbft/cometbft/abci/types"
	tmos "github.com/cometbft/cometbft/libs/os"
	dbm "github.com/cosmos/cosmos-db"
	_ "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts" // import for side-effects
	icahostkeeper "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/keeper"
	icahosttypes "github.com/cosmos/ibc-go/v8/modules/apps/27-interchain-accounts/host/types"
	ibcfee "github.com/cosmos/ibc-go/v8/modules/apps/29-fee"
	ibcfeekeeper "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/keeper" // import for side-effects
	ibcfeetypes "github.com/cosmos/ibc-go/v8/modules/apps/29-fee/types"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibcporttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	"github.com/spf13/cast"

	// "github.com/ignite/cli/ignite/pkg/cosmoscmd"

	arkeomodule "github.com/arkeonetwork/arkeo/x/arkeo"
	arkeomodulekeeper "github.com/arkeonetwork/arkeo/x/arkeo/keeper"
	arkeomoduletypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	claimmodule "github.com/arkeonetwork/arkeo/x/claim"
	claimmodulekeeper "github.com/arkeonetwork/arkeo/x/claim/keeper"
	claimmoduletypes "github.com/arkeonetwork/arkeo/x/claim/types"
	// this line is used by starport scaffolding # stargate/app/moduleImport
)

const (
	AppName = "arkeo"
)

const (
	NodeDir      = ".arkeo"
	Bech32Prefix = AccountAddressPrefix
)

var (

	// DefaultNodeHome default home directories for the application daemon
	DefaultNodeHome string
)

var (
	// Bech32PrefixAccAddr
	Bech32PrefixAccAddr  = Bech32Prefix
	Bech32PrefixAccPub   = Bech32Prefix + sdk.PrefixPublic
	Bech32PrefixValAddr  = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	Bech32PrefixValPub   = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	Bech32PrefixConsAddr = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	Bech32PrefixConsPub  = Bech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

var (

	// ModuleBasics defines the module BasicManager is in charge of setting up basic,
	// non-dependant module elements, such as codec registration
	// and genesis verification.
	ModuleBasics = module.NewBasicManager(
		auth.AppModuleBasic{},
		authzmodule.AppModuleBasic{},
		genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
		bank.AppModuleBasic{},
		capability.AppModuleBasic{},
		staking.AppModuleBasic{},
		mint.AppModuleBasic{},
		distr.AppModuleBasic{},
		gov.NewAppModuleBasic(getGovProposalHandlers()),
		params.AppModuleBasic{},
		crisis.AppModuleBasic{},
		slashing.AppModuleBasic{},
		feegrantmodule.AppModuleBasic{},
		groupmodule.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibcfee.AppModuleBasic{},
		upgrade.AppModuleBasic{},
		evidence.AppModuleBasic{},
		transfer.AppModuleBasic{},
		ica.AppModuleBasic{},
		vesting.AppModuleBasic{},
		arkeomodule.AppModuleBasic{},
		claimmodule.AppModuleBasic{},
		// this line is used by starport scaffolding # stargate/app/moduleBasic
	)

	// module account permissions
	maccPerms = map[string][]string{
		authtypes.FeeCollectorName:     nil,
		distrtypes.ModuleName:          {authtypes.Minter},
		icatypes.ModuleName:            nil,
		minttypes.ModuleName:           {authtypes.Minter},
		stakingtypes.BondedPoolName:    {authtypes.Burner, authtypes.Staking},
		stakingtypes.NotBondedPoolName: {authtypes.Burner, authtypes.Staking},
		govtypes.ModuleName:            {authtypes.Burner},
		ibctransfertypes.ModuleName:    {authtypes.Minter, authtypes.Burner},
		ibcfeetypes.ModuleName:         nil,
		claimmoduletypes.ModuleName:    {authtypes.Minter},
		arkeomoduletypes.ReserveName:   {},
		arkeomoduletypes.ModuleName:    {},
		arkeomoduletypes.ProviderName:  {},
		arkeomoduletypes.ContractName:  {},
	}
)

var (
	_ runtime.AppI            = (*ArkeoApp)(nil)
	_ servertypes.Application = (*ArkeoApp)(nil)
)

// App extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type ArkeoApp struct {
	*baseapp.BaseApp

	cdc               *codec.LegacyAmino
	appCodec          codec.Codec
	interfaceRegistry types.InterfaceRegistry
	txConfig          client.TxConfig
	invCheckPeriod    uint

	// keys to access the substores
	keys    map[string]*storetypes.KVStoreKey
	tkeys   map[string]*storetypes.TransientStoreKey
	memKeys map[string]*storetypes.MemoryStoreKey

	// keepers
	Keepers keepers.ArkeoKeepers

	// make scoped keepers public for test purposes
	ScopedIBCKeeper      capabilitykeeper.ScopedKeeper
	ScopedTransferKeeper capabilitykeeper.ScopedKeeper
	ScopedICAHostKeeper  capabilitykeeper.ScopedKeeper

	// this line is used by starport scaffolding # stargate/app/keeperDeclaration

	// mm is the module manager
	mm *module.Manager

	// sm is the simulation manager
	sm           *module.SimulationManager
	configurator module.Configurator
}

func init() {
	var err error
	DefaultNodeHome, err = clienthelpers.GetNodeHomeDirectory("." + AppName)
	if err != nil {
		panic(err)
	}
}

// this line is used by starport scaffolding # stargate/wasm/app/enabledProposals

func getGovProposalHandlers() []govclient.ProposalHandler {
	var govProposalHandlers []govclient.ProposalHandler
	// this line is used by starport scaffolding # stargate/app/govProposalHandlers

	govProposalHandlers = append(govProposalHandlers,
		paramsclient.ProposalHandler,
		// upgradeclient.LegacyProposalHandler,
		// upgradeclient.LegacyCancelProposalHandler,
		// ibcclientclient.UpdateClientProposalHandler,
		// ibcclientclient.UpgradeProposalHandler,
		// this line is used by starport scaffolding # stargate/app/govProposalHandler
	)

	return govProposalHandlers
}

// nolint
// New returns a reference to an initialized blockchain app
func NewArkeoApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	skipUpgradeHeights map[int64]bool,
	homePath string,
	invCheckPeriod uint,
	encodingConfig arekoappParams.EncodingConfig,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) *ArkeoApp {
	interfaceRegistry := encodingConfig.InterfaceRegistry

	appCodec, cdc := codec.NewProtoCodec(interfaceRegistry), codec.NewLegacyAmino()

	std.RegisterLegacyAminoCodec(cdc)
	std.RegisterInterfaces(interfaceRegistry)

	// Register IBC client state types
	ibctm.RegisterInterfaces(interfaceRegistry)

	bApp := baseapp.NewBaseApp(AppName, logger, db, encodingConfig.TxConfig.TxDecoder(), baseAppOptions...)
	bApp.SetCommitMultiStoreTracer(traceStore)
	bApp.SetVersion(version.Version)
	bApp.SetInterfaceRegistry(interfaceRegistry)
	bApp.SetTxEncoder(encodingConfig.TxConfig.TxEncoder())

	txConfig := tx.NewTxConfig(appCodec, tx.DefaultSignModes)

	keys := storetypes.NewKVStoreKeys(
		authtypes.StoreKey, banktypes.StoreKey, stakingtypes.StoreKey,
		minttypes.StoreKey, distrtypes.StoreKey, slashingtypes.StoreKey, govtypes.StoreKey,
		consensusparamtypes.StoreKey,
		paramstypes.StoreKey, ibcexported.StoreKey, upgradetypes.StoreKey, feegrant.StoreKey, evidencetypes.StoreKey,
		ibctransfertypes.StoreKey, icahosttypes.StoreKey, ibcfeetypes.StoreKey, crisistypes.StoreKey, capabilitytypes.StoreKey, group.StoreKey, authzkeeper.StoreKey,
		arkeomoduletypes.StoreKey,
		claimmoduletypes.StoreKey,
		// this line is used by starport scaffolding # stargate/app/storeKey
	)
	tkeys := storetypes.NewTransientStoreKeys(paramstypes.TStoreKey)
	memKeys := storetypes.NewMemoryStoreKeys(capabilitytypes.MemStoreKey)

	if err := bApp.RegisterStreamingServices(appOpts, keys); err != nil {
		panic(err)
	}

	app := &ArkeoApp{
		BaseApp:           bApp,
		cdc:               cdc,
		appCodec:          appCodec,
		interfaceRegistry: interfaceRegistry,
		invCheckPeriod:    invCheckPeriod,
		keys:              keys,
		tkeys:             tkeys,
		memKeys:           memKeys,
		txConfig:          txConfig,
	}

	app.Keepers.ParamsKeeper = initParamsKeeper(
		appCodec,
		cdc,
		keys[paramstypes.StoreKey],
		tkeys[paramstypes.TStoreKey],
	)
	govModuleAddr := authtypes.NewModuleAddress(govtypes.ModuleName).String()

	app.Keepers.ConsensusParamsKeeper = consensuskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[consensusparamtypes.StoreKey]),
		govModuleAddr,
		runtime.EventService{},
	)

	// set the BaseApp's parameter store
	app.SetParamStore(app.Keepers.ConsensusParamsKeeper.ParamsStore) // nolint
	// bApp.SetParamStore(app.ParamsKeeper.Subspace(baseapp.Paramspace).WithKeyTable(paramstypes.ConsensusParamsKeyTable()))

	// add capability keeper and ScopeToModule for ibc module
	app.Keepers.CapabilityKeeper = capabilitykeeper.NewKeeper(
		appCodec,
		keys[capabilitytypes.StoreKey],
		memKeys[capabilitytypes.MemStoreKey],
	)

	// grant capabilities for the ibc and ibc-transfer modules
	scopedIBCKeeper := app.Keepers.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.Keepers.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	scopedICAHostKeeper := app.Keepers.CapabilityKeeper.ScopeToModule(icahosttypes.SubModuleName)
	// this line is used by starport scaffolding # stargate/app/scopedKeeper
	app.Keepers.CapabilityKeeper.Seal()

	// add keepers
	app.Keepers.AccountKeeper = authkeeper.NewAccountKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[authtypes.StoreKey]),
		authtypes.ProtoBaseAccount,
		maccPerms,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		sdk.GetConfig().GetBech32AccountAddrPrefix(),
		govModuleAddr,
	)

	app.Keepers.AuthzKeeper = authzkeeper.NewKeeper(
		runtime.NewKVStoreService(keys[authz.ModuleName]),
		appCodec,
		app.MsgServiceRouter(), // nolint
		app.Keepers.AccountKeeper,
	)

	app.Keepers.BankKeeper = bankkeeper.NewBaseKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[banktypes.StoreKey]),
		app.Keepers.AccountKeeper,
		app.BlockedModuleAccountAddrs(),
		govModuleAddr,
		logger,
	)

	app.Keepers.StakingKeeper = stakingkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[stakingtypes.StoreKey]),
		app.Keepers.AccountKeeper,
		app.Keepers.BankKeeper,
		govModuleAddr,
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	)

	app.Keepers.MintKeeper = mintkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[minttypes.StoreKey]),
		app.Keepers.StakingKeeper,
		app.Keepers.AccountKeeper,
		app.Keepers.BankKeeper,
		authtypes.FeeCollectorName,
		govModuleAddr,
	)

	app.Keepers.DistrKeeper = distrkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[distrtypes.StoreKey]),
		app.Keepers.AccountKeeper,
		app.Keepers.BankKeeper,
		app.Keepers.StakingKeeper,
		authtypes.FeeCollectorName,
		govModuleAddr,
	)

	app.Keepers.SlashingKeeper = slashingkeeper.NewKeeper(
		appCodec,
		cdc,
		runtime.NewKVStoreService(keys[slashingtypes.StoreKey]),
		app.Keepers.StakingKeeper,
		govModuleAddr,
	)

	app.Keepers.CrisisKeeper = *crisiskeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[crisistypes.StoreKey]),
		invCheckPeriod,
		app.Keepers.BankKeeper,
		authtypes.FeeCollectorName,
		govModuleAddr,
		app.Keepers.AccountKeeper.AddressCodec(),
	)

	app.Keepers.FeeGrantKeeper = feegrantkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[feegrant.StoreKey]),
		app.Keepers.AccountKeeper,
	)

	app.Keepers.UpgradeKeeper = upgradekeeper.NewKeeper(
		skipUpgradeHeights,
		runtime.NewKVStoreService(keys[upgradetypes.StoreKey]),
		appCodec,
		homePath,
		app.BaseApp,
		govModuleAddr,
	)
	app.Keepers.UpgradeKeeper.SetUpgradeHandler("claim-fix-v1.0.10", func(ctx context.Context, plan upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
		app.Keepers.ClaimKeeper.RunUpgradeClaimFix(ctx)
		return app.mm.RunMigrations(ctx, app.configurator, fromVM)
	})

	groupConfig := group.DefaultConfig()
	/*
		Example of setting group params:
		groupConfig.MaxMetadataLen = 1000
	*/
	app.Keepers.GroupKeeper = groupkeeper.NewKeeper(
		keys[group.StoreKey],
		appCodec,
		app.MsgServiceRouter(), // nolint
		app.Keepers.AccountKeeper,
		groupConfig,
	)

	app.Keepers.ClaimKeeper = claimmodulekeeper.NewKeeper(
		appCodec,
		keys[claimmoduletypes.StoreKey],
		app.Keepers.AccountKeeper,
		app.Keepers.BankKeeper,
		keys[claimmoduletypes.MemStoreKey],
		app.GetSubspace(claimmoduletypes.ModuleName),
		logger,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks

	app.Keepers.StakingKeeper.SetHooks(stakingtypes.NewMultiStakingHooks(app.Keepers.DistrKeeper.Hooks(), app.Keepers.SlashingKeeper.Hooks(), app.Keepers.ClaimKeeper.Hooks()))

	// ... other modules keepers

	govRouter := govv1beta1.NewRouter()
	govRouter.AddRoute(govtypes.RouterKey, govv1beta1.ProposalHandler).AddRoute(
		paramproposal.RouterKey, params.NewParamChangeProposalHandler(app.Keepers.ParamsKeeper))

	// Create IBC Keeper
	app.Keepers.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		keys[ibcexported.StoreKey],
		app.GetSubspace(ibcexported.ModuleName),
		app.Keepers.StakingKeeper,
		app.Keepers.UpgradeKeeper,
		scopedIBCKeeper,
		govModuleAddr,
	)

	app.Keepers.IBCFeeKeeper = ibcfeekeeper.NewKeeper(
		appCodec,
		keys[ibcfeetypes.StoreKey],
		app.Keepers.IBCKeeper.ChannelKeeper,
		app.Keepers.IBCKeeper.ChannelKeeper,
		app.Keepers.IBCKeeper.PortKeeper, 
		app.Keepers.AccountKeeper, 
		app.Keepers.BankKeeper,
	)

	// Create Transfer Keepers
	app.Keepers.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		keys[ibctransfertypes.StoreKey],
		app.GetSubspace(ibctransfertypes.ModuleName),
		app.Keepers.IBCFeeKeeper,
		app.Keepers.IBCKeeper.ChannelKeeper,
		app.Keepers.IBCKeeper.PortKeeper,
		app.Keepers.AccountKeeper,
		app.Keepers.BankKeeper,
		scopedTransferKeeper,
		govModuleAddr,
	)
	transferModule := transfer.NewAppModule(app.Keepers.TransferKeeper)
	transferIBCModule := transfer.NewIBCModule(app.Keepers.TransferKeeper)

	app.Keepers.ICAHostKeeper = icahostkeeper.NewKeeper(
		appCodec,
		keys[icahosttypes.StoreKey],
		app.GetSubspace(icahosttypes.SubModuleName),
		app.Keepers.IBCFeeKeeper,
		app.Keepers.IBCKeeper.ChannelKeeper,
		app.Keepers.IBCKeeper.PortKeeper,
		app.Keepers.AccountKeeper,
		scopedICAHostKeeper,
		app.MsgServiceRouter(), // nolint
		govModuleAddr,
	)
	icaModule := ica.NewAppModule(nil, &app.Keepers.ICAHostKeeper)
	icaHostIBCModule := icahost.NewIBCModule(app.Keepers.ICAHostKeeper)

	app.Keepers.ICAHostKeeper.WithQueryRouter(app.GRPCQueryRouter()) // nolint

	// Create evidence Keeper for to register the IBC light client misbehavior evidence route
	evidenceKeeper := evidencekeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[evidencetypes.StoreKey]),
		app.Keepers.StakingKeeper,
		app.Keepers.SlashingKeeper,
		app.Keepers.AccountKeeper.AddressCodec(),
		runtime.ProvideCometInfoService(),
	)
	// If evidence needs to be handled for the app, set routes in router here and seal
	app.Keepers.EvidenceKeeper = *evidenceKeeper

	// this line is used by starport scaffolding # stargate/app/keeperDefinition

	// Create static IBC router, add transfer route, then set and seal it
	ibcRouter := ibcporttypes.NewRouter()
	ibcRouter.AddRoute(icahosttypes.SubModuleName, icaHostIBCModule)
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferIBCModule)

	// this line is used by starport scaffolding # ibc/app/router
	app.Keepers.IBCKeeper.SetRouter(ibcRouter)

	govConfig := govtypes.DefaultConfig()
	app.Keepers.GovKeeper = govkeeper.NewKeeper(
		appCodec,
		runtime.NewKVStoreService(keys[govtypes.StoreKey]),
		app.Keepers.AccountKeeper,
		app.Keepers.BankKeeper,
		app.Keepers.StakingKeeper,
		app.Keepers.DistrKeeper,
		app.MsgServiceRouter(),
		govConfig,
		govModuleAddr,
	)

	app.Keepers.GovKeeper.SetLegacyRouter(govRouter)

	app.Keepers.GovKeeper.SetHooks(
		govtypes.NewMultiGovHooks(
			app.Keepers.ClaimKeeper.Hooks(),
		),
	)

	app.Keepers.ArkeoKeeper = *arkeomodulekeeper.NewKVStore(
		appCodec,
		keys[arkeomoduletypes.StoreKey],
		keys[arkeomoduletypes.MemStoreKey],
		app.GetSubspace(arkeomoduletypes.ModuleName),
		app.Keepers.BankKeeper,
		app.Keepers.AccountKeeper,
		*app.Keepers.StakingKeeper,
		govModuleAddr,
		logger,
		app.Keepers.MintKeeper,
		app.Keepers.DistrKeeper,
	)

	/****  Module Options ****/

	// NOTE: we may consider parsing `appOpts` inside module constructors. For the moment
	// we prefer to be more strict in what arguments the modules expect.
	skipGenesisInvariants := cast.ToBool(appOpts.Get(crisis.FlagSkipGenesisInvariants))

	// NOTE: Any module instantiated in the module manager that is later modified
	// must be passed by reference here.

	app.mm = module.NewManager(
		genutil.NewAppModule(
			app.Keepers.AccountKeeper, app.Keepers.StakingKeeper, app.BaseApp,
			encodingConfig.TxConfig,
		),
		auth.NewAppModule(appCodec, app.Keepers.AccountKeeper, nil, app.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.Keepers.AuthzKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.interfaceRegistry),
		vesting.NewAppModule(app.Keepers.AccountKeeper, app.Keepers.BankKeeper),
		bank.NewAppModule(appCodec, app.Keepers.BankKeeper, app.Keepers.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.Keepers.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.FeeGrantKeeper, app.interfaceRegistry),
		groupmodule.NewAppModule(appCodec, app.Keepers.GroupKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.interfaceRegistry),
		crisis.NewAppModule(&app.Keepers.CrisisKeeper, skipGenesisInvariants, app.GetSubspace(crisistypes.ModuleName)),
		gov.NewAppModule(appCodec, app.Keepers.GovKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.Keepers.MintKeeper, app.Keepers.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.Keepers.SlashingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		distr.NewAppModule(appCodec, app.Keepers.DistrKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		staking.NewAppModule(appCodec, app.Keepers.StakingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		upgrade.NewAppModule(app.Keepers.UpgradeKeeper, app.Keepers.AccountKeeper.AddressCodec()),
		evidence.NewAppModule(app.Keepers.EvidenceKeeper),
		ibc.NewAppModule(app.Keepers.IBCKeeper),
		ibcfee.NewAppModule(app.Keepers.IBCFeeKeeper),
		params.NewAppModule(app.Keepers.ParamsKeeper),
		transferModule,
		icaModule,
		arkeomodule.NewAppModule(appCodec, app.Keepers.ArkeoKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, *app.Keepers.StakingKeeper),
		claimmodule.NewAppModule(appCodec, app.Keepers.ClaimKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper),
		// this line is used by starport scaffolding # stargate/app/appModule
	)

	// During begin block slashing happens after distr.BeginBlocker so that
	// there is nothing left over in the validator fee pool, so as to keep the
	// CanWithdrawInvariant invariant.
	// NOTE: staking module is required if HistoricalEntries param > 0
	app.mm.SetOrderBeginBlockers(
		govtypes.ModuleName,
		// upgrades should be run first
		upgradetypes.ModuleName,
		capabilitytypes.ModuleName,
		minttypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		evidencetypes.ModuleName,
		stakingtypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		crisistypes.ModuleName,
		ibcfeetypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		genutiltypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		vestingtypes.ModuleName,
		arkeomoduletypes.ModuleName,
		claimmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/beginBlockers
	)

	app.mm.SetOrderEndBlockers(
		crisistypes.ModuleName,
		govtypes.ModuleName,
		stakingtypes.ModuleName,
		ibcfeetypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		slashingtypes.ModuleName,
		minttypes.ModuleName,
		genutiltypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		arkeomoduletypes.ModuleName,
		claimmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/endBlockers
	)

	// NOTE: The genutils module must occur after staking so that pools are
	// properly initialized with tokens from genesis accounts.
	// NOTE: Capability module must occur first so that it can initialize any capabilities
	// so that other modules that want to create or claim capabilities afterwards in InitChain
	// can do so safely.
	app.mm.SetOrderInitGenesis(
		capabilitytypes.ModuleName,
		authtypes.ModuleName,
		banktypes.ModuleName,
		distrtypes.ModuleName,
		stakingtypes.ModuleName,
		slashingtypes.ModuleName,
		govtypes.ModuleName,
		minttypes.ModuleName,
		crisistypes.ModuleName,
		genutiltypes.ModuleName,
		ibcfeetypes.ModuleName,
		ibctransfertypes.ModuleName,
		ibcexported.ModuleName,
		icatypes.ModuleName,
		evidencetypes.ModuleName,
		authz.ModuleName,
		feegrant.ModuleName,
		group.ModuleName,
		paramstypes.ModuleName,
		upgradetypes.ModuleName,
		vestingtypes.ModuleName,
		arkeomoduletypes.ModuleName,
		claimmoduletypes.ModuleName,
		// this line is used by starport scaffolding # stargate/app/initGenesis
	)

	// Uncomment if you want to set a custom migration order here.
	// app.mm.SetOrderMigrations(custom order)

	app.mm.RegisterInvariants(&app.Keepers.CrisisKeeper)

	app.configurator = module.NewConfigurator(app.appCodec, app.MsgServiceRouter(), app.GRPCQueryRouter())
	err := app.mm.RegisterServices(app.configurator)
	if err != nil {
		panic(fmt.Errorf("failed to register services: %s", err))
	}

	app.RegisterUpgradeHandlers()
	autocliv1.RegisterQueryServer(app.GRPCQueryRouter(), runtimeservices.NewAutoCLIQueryService(app.mm.Modules))
	reflectionSvc, err := runtimeservices.NewReflectionService()
	if err != nil {
		panic(err)
	}
	reflectionv1.RegisterReflectionServiceServer(app.GRPCQueryRouter(), reflectionSvc)

	protoFiles, err := proto.MergedRegistry()
	if err != nil {
		panic(err)
	}

	err = msgservice.ValidateProtoAnnotations(protoFiles)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	// create the simulation manager and define the order of the modules for deterministic simulations
	app.sm = module.NewSimulationManager(
		auth.NewAppModule(appCodec, app.Keepers.AccountKeeper, authsims.RandomGenesisAccounts, app.GetSubspace(authtypes.ModuleName)),
		authzmodule.NewAppModule(appCodec, app.Keepers.AuthzKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.interfaceRegistry),
		bank.NewAppModule(appCodec, app.Keepers.BankKeeper, app.Keepers.AccountKeeper, app.GetSubspace(banktypes.ModuleName)),
		capability.NewAppModule(appCodec, *app.Keepers.CapabilityKeeper, false),
		feegrantmodule.NewAppModule(appCodec, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.FeeGrantKeeper, app.interfaceRegistry),
		gov.NewAppModule(appCodec, app.Keepers.GovKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.GetSubspace(govtypes.ModuleName)),
		mint.NewAppModule(appCodec, app.Keepers.MintKeeper, app.Keepers.AccountKeeper, nil, app.GetSubspace(minttypes.ModuleName)),
		staking.NewAppModule(appCodec, app.Keepers.StakingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.GetSubspace(stakingtypes.ModuleName)),
		distr.NewAppModule(appCodec, app.Keepers.DistrKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.StakingKeeper, app.GetSubspace(distrtypes.ModuleName)),
		slashing.NewAppModule(appCodec, app.Keepers.SlashingKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.Keepers.StakingKeeper, app.GetSubspace(slashingtypes.ModuleName), app.interfaceRegistry),
		params.NewAppModule(app.Keepers.ParamsKeeper),
		groupmodule.NewAppModule(appCodec, app.Keepers.GroupKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, app.interfaceRegistry),
		evidence.NewAppModule(app.Keepers.EvidenceKeeper),
		ibc.NewAppModule(app.Keepers.IBCKeeper),
		transferModule,
		arkeomodule.NewAppModule(appCodec, app.Keepers.ArkeoKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper, *app.Keepers.StakingKeeper),
		claimmodule.NewAppModule(appCodec, app.Keepers.ClaimKeeper, app.Keepers.AccountKeeper, app.Keepers.BankKeeper),
		// this line is used by starport scaffolding # stargate/app/appModule
	)
	app.sm.RegisterStoreDecoders()

	// initialize stores
	app.MountKVStores(keys)
	app.MountTransientStores(tkeys)
	app.MountMemoryStores(memKeys)

	// initialize BaseApp
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)

	anteHandler, err := ante.NewAnteHandler(
		ante.HandlerOptions{
			AccountKeeper:   app.Keepers.AccountKeeper,
			BankKeeper:      app.Keepers.BankKeeper,
			SignModeHandler: encodingConfig.TxConfig.SignModeHandler(),
			FeegrantKeeper:  app.Keepers.FeeGrantKeeper,
			SigGasConsumer:  ante.DefaultSigVerificationGasConsumer,
		},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create AnteHandler: %s", err))
	}

	postHandler, err := posthandler.NewPostHandler(
		posthandler.HandlerOptions{},
	)
	if err != nil {
		panic(fmt.Errorf("failed to create PostHandler: %s", err))
	}

	app.SetAnteHandler(anteHandler)
	app.SetInitChainer(app.InitChainer)
	app.SetPreBlocker(app.PreBlocker)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)
	app.SetPostHandler(postHandler)

	if loadLatest {
		if err := app.LoadLatestVersion(); err != nil {
			tmos.Exit(err.Error())
		}
	}

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper
	// this line is used by starport scaffolding # stargate/app/beforeInitReturn

	return app
}

func (app *ArkeoApp) PreBlocker(ctx sdk.Context, _ *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	return app.mm.PreBlock(ctx)
}

// Name returns the name of the App
func (app *ArkeoApp) Name() string { return app.BaseApp.Name() }

// GetBaseApp returns the base app of the application
func (app ArkeoApp) GetBaseApp() *baseapp.BaseApp { return app.BaseApp }

// BeginBlocker application updates every begin block
func (app *ArkeoApp) BeginBlocker(ctx sdk.Context) (sdk.BeginBlock, error) {
	return app.mm.BeginBlock(ctx)
}

// EndBlocker application updates every end block
func (app *ArkeoApp) EndBlocker(ctx sdk.Context) (sdk.EndBlock, error) {
	return app.mm.EndBlock(ctx)
}

// InitChainer application update at chain initialization
func (app *ArkeoApp) InitChainer(ctx sdk.Context, req *abci.RequestInitChain) (*abci.ResponseInitChain, error) {
	var genesisState GenesisState
	app.cdc.MustUnmarshalJSON(req.AppStateBytes, &genesisState)
	err := app.Keepers.UpgradeKeeper.SetModuleVersionMap(ctx, app.mm.GetVersionMap())
	if err != nil {
		panic(err)
	}
	mintGen := minttypes.GenesisState{
		Minter: minttypes.Minter{
			Inflation:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
			AnnualProvisions: math.LegacyMustNewDecFromStr("0.000000000000000000"),
		},
		Params: minttypes.Params{
			MintDenom:           "uarkeo",
			InflationRateChange: math.LegacyMustNewDecFromStr("0.000000000000000000"),
			InflationMax:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
			InflationMin:        math.LegacyMustNewDecFromStr("0.000000000000000000"),
			GoalBonded:          math.LegacyNewDec(670000000000000000),
			BlocksPerYear:       6311520,
		},
	}
	genesisState[minttypes.ModuleName] = app.cdc.MustMarshalJSON(&mintGen)
	return app.mm.InitGenesis(ctx, app.appCodec, genesisState)
}

// nolint
// LoadHeight loads a particular height
func (app *ArkeoApp) LoadHeight(height int64) error {
	return app.LoadVersion(height)
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *ArkeoApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[authtypes.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

// BlockedModuleAccountAddrs returns all the app's blocked module account
// addresses.
func (app *ArkeoApp) BlockedModuleAccountAddrs() map[string]bool {
	modAccAddrs := app.ModuleAccountAddrs()
	delete(modAccAddrs, authtypes.NewModuleAddress(govtypes.ModuleName).String())

	return modAccAddrs
}

// LegacyAmino returns SimApp's amino codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ArkeoApp) LegacyAmino() *codec.LegacyAmino {
	return app.cdc
}

// AppCodec returns an app codec.
//
// NOTE: This is solely to be used for testing purposes as it may be desirable
// for modules to register their own custom testing types.
func (app *ArkeoApp) AppCodec() codec.Codec {
	return app.appCodec
}

// InterfaceRegistry returns an InterfaceRegistry
func (app *ArkeoApp) InterfaceRegistry() types.InterfaceRegistry {
	return app.interfaceRegistry
}

// GetKey returns the KVStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ArkeoApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	return app.keys[storeKey]
}

// GetTKey returns the TransientStoreKey for the provided store key.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ArkeoApp) GetTKey(storeKey string) *storetypes.TransientStoreKey {
	return app.tkeys[storeKey]
}

// GetMemKey returns the MemStoreKey for the provided mem key.
//
// NOTE: This is solely used for testing purposes.
func (app *ArkeoApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	return app.memKeys[storeKey]
}

// GetSubspace returns a param subspace for a given module name.
//
// NOTE: This is solely to be used for testing purposes.
func (app *ArkeoApp) GetSubspace(moduleName string) paramstypes.Subspace {
	subspace, _ := app.Keepers.ParamsKeeper.GetSubspace(moduleName)
	return subspace
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *ArkeoApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	clientCtx := apiSvr.ClientCtx
	// Register new tx routes from grpc-gateway.
	authtx.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)
	// Register new cometbft queries routes from grpc-gateway.
	cmtservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register node gRPC service for grpc-gateway
	nodeservice.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// Register grpc-gateway routes for all modules.
	ModuleBasics.RegisterGRPCGatewayRoutes(clientCtx, apiSvr.GRPCGatewayRouter)

	// register app's OpenAPI routes.
	if apiConfig.Swagger {
		if err := RegisterSwaggerAPI(apiSvr); err != nil {
			panic(err)
		}
	}
}

// RegisterTxService implements the Application.RegisterTxService method.
func (app *ArkeoApp) RegisterTxService(clientCtx client.Context) {
	authtx.RegisterTxService(app.BaseApp.GRPCQueryRouter(), clientCtx, app.BaseApp.Simulate, app.interfaceRegistry)
}

// nolint
// RegisterTendermintService implements the Application.RegisterTendermintService method.
func (app *ArkeoApp) RegisterTendermintService(clientCtx client.Context) {
	cmtservice.RegisterTendermintService(
		clientCtx,
		app.BaseApp.GRPCQueryRouter(),
		app.interfaceRegistry,
		app.Query,
	)
}

// nolint
func (app *ArkeoApp) RegisterNodeService(clientCtx client.Context, cfg config.Config) {
	nodeservice.RegisterNodeService(clientCtx, app.GRPCQueryRouter(), cfg)
}

// GetMaccPerms returns a copy of the module account permissions
func GetMaccPerms() map[string][]string {
	dupMaccPerms := make(map[string][]string)
	for k, v := range maccPerms {
		dupMaccPerms[k] = v
	}
	return dupMaccPerms
}

// initParamsKeeper init params keeper and its subspaces
func initParamsKeeper(appCodec codec.BinaryCodec, legacyAmino *codec.LegacyAmino, key, tkey storetypes.StoreKey) paramskeeper.Keeper {
	paramsKeeper := paramskeeper.NewKeeper(appCodec, legacyAmino, key, tkey)

	paramsKeeper.Subspace(authtypes.ModuleName)
	paramsKeeper.Subspace(banktypes.ModuleName)
	paramsKeeper.Subspace(stakingtypes.ModuleName)
	paramsKeeper.Subspace(minttypes.ModuleName)
	paramsKeeper.Subspace(distrtypes.ModuleName)
	paramsKeeper.Subspace(slashingtypes.ModuleName)
	paramsKeeper.Subspace(govtypes.ModuleName)
	paramsKeeper.Subspace(crisistypes.ModuleName)
	paramsKeeper.Subspace(ibcfeetypes.ModuleName)
	paramsKeeper.Subspace(ibctransfertypes.ModuleName) //nolint:staticcheck
	paramsKeeper.Subspace(ibcexported.ModuleName)
	paramsKeeper.Subspace(icahosttypes.SubModuleName)
	paramsKeeper.Subspace(arkeomoduletypes.ModuleName)
	paramsKeeper.Subspace(arkeomoduletypes.ProviderName)
	paramsKeeper.Subspace(arkeomoduletypes.ContractName)
	paramsKeeper.Subspace(claimmoduletypes.ModuleName)
	// this line is used by starport scaffolding # stargate/app/paramSubspace

	return paramsKeeper
}

// SimulationManager implements the SimulationApp interface
func (app *ArkeoApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

func (app *ArkeoApp) AutoCliOpts() autocli.AppOptions {
	modules := make(map[string]appmodule.AppModule, 0)

	for _, m := range app.mm.Modules {
		if moduleWithName, ok := m.(module.HasName); ok {
			moduleName := moduleWithName.Name()
			if appModule, ok := moduleWithName.(appmodule.AppModule); ok {
				modules[moduleName] = appModule
			}
		}
	}

	return autocli.AppOptions{
		Modules:               modules,
		ModuleOptions:         runtimeservice.ExtractAutoCLIOptions(app.mm.Modules),
		AddressCodec:          authcodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix()),
		ValidatorAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ValidatorAddrPrefix()),
		ConsensusAddressCodec: authcodec.NewBech32Codec(sdk.GetConfig().GetBech32ConsensusAddrPrefix()),
	}
}
