package types

import (
	fmt "fmt"
	time "time"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	KeyClaimDenom            = []byte("ClaimDenom")
	DefaultClaimDenom string = "uarkeo"
)

var (
	KeyDurationUntilDecay                   = []byte("DurationUntilDecay")
	DefaultDurationUntilDecay time.Duration = time.Hour
)

var (
	KeyDurationOfDecay                   = []byte("DurationOfDecay")
	DefaultDurationOfDecay time.Duration = time.Hour
)

var (
	KeyAirdropStartTime               = []byte("AirdropStartTime")
	DeafultAirdropStartTime time.Time = time.Now().UTC()
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(claimDenom string, airdropStartTime time.Time, durationUntilDecay, durationOfDecay time.Duration) Params {
	return Params{
		ClaimDenom:         claimDenom,
		AirdropStartTime:   airdropStartTime,
		DurationUntilDecay: durationUntilDecay,
		DurationOfDecay:    durationOfDecay,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		ClaimDenom:         DefaultClaimDenom,
		DurationUntilDecay: DefaultDurationUntilDecay,
		DurationOfDecay:    DefaultDurationOfDecay,
		AirdropStartTime:   DeafultAirdropStartTime,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyAirdropStartTime, &p.AirdropStartTime, validateAirdropStartTime),
		paramtypes.NewParamSetPair(KeyDurationUntilDecay, &p.DurationUntilDecay, validateDurationUntilDecay),
		paramtypes.NewParamSetPair(KeyDurationOfDecay, &p.DurationOfDecay, validateDurationOfDecay),
		paramtypes.NewParamSetPair(KeyClaimDenom, &p.ClaimDenom, validateClaimDenom),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

func validateAirdropStartTime(i interface{}) error {
	_, ok := i.(time.Time)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateDurationOfDecay(i interface{}) error {
	_, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateDurationUntilDecay(i interface{}) error {
	_, ok := i.(time.Duration)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateClaimDenom(i interface{}) error {
	_, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}
