package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// Hooks wrapper struct
type Hooks struct {
	k Keeper
}

// Return the wrapper struct
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

var _ govtypes.GovHooks = Hooks{}
var _ stakingtypes.StakingHooks = Hooks{}

// governance hooks
func (h Hooks) AfterProposalSubmission(ctx sdk.Context, proposalID uint64) {}
func (h Hooks) AfterProposalDeposit(ctx sdk.Context, proposalID uint64, depositorAddr sdk.AccAddress) {
}

func (h Hooks) AfterProposalVote(ctx sdk.Context, proposalID uint64, voterAddr sdk.AccAddress) {
	h.k.AfterProposalVote(ctx, proposalID, voterAddr)
}
func (h Hooks) AfterProposalFailedMinDeposit(ctx sdk.Context, proposalID uint64)  {}
func (h Hooks) AfterProposalVotingPeriodEnded(ctx sdk.Context, proposalID uint64) {}

// staking hooks
func (h Hooks) AfterValidatorCreated(ctx sdk.Context, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeValidatorModified(ctx sdk.Context, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorRemoved(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}

func (h Hooks) AfterValidatorBonded(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterValidatorBeginUnbonding(ctx sdk.Context, consAddr sdk.ConsAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationCreated(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationSharesModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) BeforeDelegationRemoved(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	return nil
}
func (h Hooks) AfterDelegationModified(ctx sdk.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	_ = h.k.AfterDelegationModified(ctx, delAddr, valAddr)
	return nil
}
func (h Hooks) BeforeValidatorSlashed(ctx sdk.Context, valAddr sdk.ValAddress, fraction sdk.Dec) error {
	return nil
}
