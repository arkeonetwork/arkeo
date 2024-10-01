package types

import (
	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams() Params {
	return Params{
		CommunityPoolPercentage:   math.LegacyMustNewDecFromStr("0.100000000000000000"),
		DevFundPercentage:         math.LegacyMustNewDecFromStr("0.200000000000000000"),
		GrantFundPercentage:       math.LegacyMustNewDecFromStr("0.000000000000000000"),
		InflationChangePercentage: math.LegacyMustNewDecFromStr("0.030000000000000000"),
		InflationMin:              math.LegacyMustNewDecFromStr("0.020000000000000000"),
		InflationMax:              math.LegacyMustNewDecFromStr("0.050000000000000000"),
		GoalBonded:                math.LegacyMustNewDecFromStr("0.670000000000000000"),
		BlockPerYear:              6311520,
		EmissionCurve:             6,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams()
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}
