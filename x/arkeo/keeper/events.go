package keeper

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) EmitBondProviderEvent(ctx cosmos.Context, bond cosmos.Int, msg *types.MsgBondProvider) error {
	return ctx.EventManager().EmitTypedEvent(
		&types.EventBondProvider{
			Provider: msg.Provider,
			Creator:  msg.Creator,
			Service:  msg.Service,
			Bond:     msg.Bond,
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

func (k msgServer) EmitModProviderEvent(ctx cosmos.Context, msg *types.MsgModProvider, provider *types.Provider) error {
	return ctx.EventManager().EmitTypedEvent(
		&types.EventModProvider{
			Creator:             msg.Creator,
			Provider:            provider.PubKey,
			Service:             provider.Service.String(),
			MetadataURI:         provider.MetadataUri,
			MetadataNonce:       provider.MetadataNonce,
			Status:              types.ProviderStatus(provider.Status),
			MinContractDuration: provider.MinContractDuration,
			MaxContractDuration: provider.MaxContractDuration,
			SubscriptionRate:    provider.SubscriptionRate,
			PayAsYouGoRate:      provider.PayAsYouGoRate,
			Bond:                provider.Bond,
			SettlementDuration:  provider.SettlementDuration,
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

func (mgr Manager) ContractSettlementEvent(ctx cosmos.Context, debt, valIncome cosmos.Int, contract *types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			types.NewContractSettlementEvent(debt, valIncome, contract),
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
