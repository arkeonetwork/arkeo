package keeper

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BondProviderEvent(ctx cosmos.Context, bond cosmos.Int, msg *types.MsgBondProvider) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			types.NewBondProviderEvent(bond, msg),
		},
	)
}

func (k msgServer) CloseContractEvent(ctx cosmos.Context, contract *types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			types.NewCloseContractEvent(contract),
		},
	)
}

func (k msgServer) ModProviderEvent(ctx cosmos.Context, provider *types.Provider) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			types.NewModProviderEvent(provider),
		},
	)
}

func (k msgServer) OpenContractEvent(ctx cosmos.Context, openCost int64, contract *types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			types.NewOpenContractEvent(openCost, contract),
		},
	)
}

func (mgr Manager) ContractSettlementEvent(ctx cosmos.Context, debt, valIncome cosmos.Int, contract *types.Contract, nonce int64) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			types.NewContractSettlementEvent(debt, valIncome, contract, nonce),
		},
	)
}

func (mgr Manager) ValidatorPayoutEvent(ctx cosmos.Context, acc cosmos.AccAddress, rwd cosmos.Int) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			types.NewValidatorPayoutEvent(acc, rwd),
		},
	)
}
