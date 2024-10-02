package app

import (
	"encoding/json"

	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/codec"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState(cdc codec.JSONCodec) GenesisState {
	defaultGenesis := ModuleBasics.DefaultGenesis(cdc)
	// set mint module params for genesis state
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
			BlocksPerYear:       5256666,
		},
	}
	defaultGenesis[minttypes.ModuleName] = cdc.MustMarshalJSON(&mintGen)

	return defaultGenesis
}
