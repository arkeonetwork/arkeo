package arkeo

import (
	"math/rand"

	arkeosimulation "github.com/arkeonetwork/arkeo/x/arkeo/simulation"

	"github.com/arkeonetwork/arkeo/testutil/sample"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// avoid unused import issue
var (
	_ = sample.AccAddress
	_ = arkeosimulation.FindAccount
	// _ = simappparams.StakePerAccount
	_ = simulation.MsgEntryKind
	_ = baseapp.Paramspace
)

const (
	opWeightMsgBondProvider = "op_weight_msg_bond_provider" // nolint
	// TODO: Determine the simulation weight value
	defaultWeightMsgBondProvider int = 100

	opWeightMsgModProvider = "op_weight_msg_mod_provider" // nolint
	// TODO: Determine the simulation weight value
	defaultWeightMsgModProvider int = 100

	opWeightMsgOpenContract = "op_weight_msg_open_contract" // nolint
	// TODO: Determine the simulation weight value
	defaultWeightMsgOpenContract int = 100

	opWeightMsgCloseContract = "op_weight_msg_close_contract" // nolint
	// TODO: Determine the simulation weight value
	defaultWeightMsgCloseContract int = 100

	opWeightMsgClaimContractIncome = "op_weight_msg_claim_contract_income" // nolint
	// TODO: Determine the simulation weight value
	defaultWeightMsgClaimContractIncome int = 100

	opWeightMsgSetVersion = "op_weight_msg_set_version" // nolint
	// TODO: Determine the simulation weight value
	defaultWeightMsgSetVersion int = 100

	// this line is used by starport scaffolding # simapp/module/const
)

// GenerateGenesisState creates a randomized GenState of the module
func (AppModule) GenerateGenesisState(simState *module.SimulationState) {
	accs := make([]string, len(simState.Accounts))
	for i, acc := range simState.Accounts {
		accs[i] = acc.Address.String()
	}
	arkeoGenesis := types.GenesisState{
		Params: types.DefaultParams(),
		// this line is used by starport scaffolding # simapp/module/genesisState
	}
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(&arkeoGenesis)
}

// ProposalContents doesn't return any content functions for governance proposals
func (AppModule) ProposalContents(_ module.SimulationState) []simtypes.WeightedProposalContent { //nolint:staticcheck
	return nil
}

// RandomizedParams creates randomized  param changes for the simulator
func (am AppModule) RandomizedParams(_ *rand.Rand) []simtypes.LegacyParamChange {
	return []simtypes.LegacyParamChange{}
}

// RegisterStoreDecoder registers a decoder
func (am AppModule) RegisterStoreDecoder(_ simtypes.StoreDecoderRegistry) {}

// WeightedOperations returns the all the gov module operations with their respective weights.
func (am AppModule) WeightedOperations(simState module.SimulationState) []simtypes.WeightedOperation {
	operations := make([]simtypes.WeightedOperation, 0)

	var weightMsgBondProvider int
	simState.AppParams.GetOrGenerate(opWeightMsgBondProvider, &weightMsgBondProvider, nil,
		func(_ *rand.Rand) {
			weightMsgBondProvider = defaultWeightMsgBondProvider
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgBondProvider,
		arkeosimulation.SimulateMsgBondProvider(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgModProvider int
	simState.AppParams.GetOrGenerate(opWeightMsgModProvider, &weightMsgModProvider, nil,
		func(_ *rand.Rand) {
			weightMsgModProvider = defaultWeightMsgModProvider
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgModProvider,
		arkeosimulation.SimulateMsgModProvider(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgOpenContract int
	simState.AppParams.GetOrGenerate(opWeightMsgOpenContract, &weightMsgOpenContract, nil,
		func(_ *rand.Rand) {
			weightMsgOpenContract = defaultWeightMsgOpenContract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgOpenContract,
		arkeosimulation.SimulateMsgOpenContract(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgCloseContract int
	simState.AppParams.GetOrGenerate(opWeightMsgCloseContract, &weightMsgCloseContract, nil,
		func(_ *rand.Rand) {
			weightMsgCloseContract = defaultWeightMsgCloseContract
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgCloseContract,
		arkeosimulation.SimulateMsgCloseContract(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgClaimContractIncome int
	simState.AppParams.GetOrGenerate(opWeightMsgClaimContractIncome, &weightMsgClaimContractIncome, nil,
		func(_ *rand.Rand) {
			weightMsgClaimContractIncome = defaultWeightMsgClaimContractIncome
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgClaimContractIncome,
		arkeosimulation.SimulateMsgClaimContractIncome(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	var weightMsgSetVersion int
	simState.AppParams.GetOrGenerate(opWeightMsgSetVersion, &weightMsgSetVersion, nil,
		func(_ *rand.Rand) {
			weightMsgSetVersion = defaultWeightMsgSetVersion
		},
	)
	operations = append(operations, simulation.NewWeightedOperation(
		weightMsgSetVersion,
		arkeosimulation.SimulateMsgSetVersion(am.accountKeeper, am.bankKeeper, am.keeper),
	))

	// this line is used by starport scaffolding # simapp/module/operation

	return operations
}
