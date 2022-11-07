package keeper

import (
	"strconv"

	"mercury/common/cosmos"
	"mercury/x/mercury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BondProviderEvent(ctx cosmos.Context, bond cosmos.Int, msg *types.MsgBondProvider) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeProviderBond,
				sdk.NewAttribute("pubkey", msg.PubKey.String()),
				sdk.NewAttribute("chain", msg.Chain),
				sdk.NewAttribute("bond_rel", msg.Bond.String()),
				sdk.NewAttribute("bond_abs", bond.String()),
			),
		},
	)
}

func (k msgServer) CloseContractEvent(ctx cosmos.Context, msg *types.MsgCloseContract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeCloseContract,
				sdk.NewAttribute("pubkey", msg.PubKey.String()),
				sdk.NewAttribute("chain", msg.Chain),
				sdk.NewAttribute("client", msg.Client.String()),
				sdk.NewAttribute("delegate", msg.Delegate.String()),
			),
		},
	)
}

func (k msgServer) ModProviderEvent(ctx cosmos.Context, provider types.Provider) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeProviderMod,
				sdk.NewAttribute("pubkey", provider.PubKey.String()),
				sdk.NewAttribute("chain", provider.Chain.String()),
				sdk.NewAttribute("metadata_uri", provider.MetadataURI),
				sdk.NewAttribute("metadata_nonce", strconv.FormatUint(provider.MetadataNonce, 10)),
				sdk.NewAttribute("status", provider.Status.String()),
				sdk.NewAttribute("min_contract_duration", strconv.FormatInt(provider.MinContractDuration, 10)),
				sdk.NewAttribute("max_contract_duration", strconv.FormatInt(provider.MaxContractDuration, 10)),
				sdk.NewAttribute("subscription_rate", strconv.FormatInt(provider.SubscriptionRate, 10)),
				sdk.NewAttribute("pay-as-you-go_rate", strconv.FormatInt(provider.PayAsYouGoRate, 10)),
			),
		},
	)
}

func (k msgServer) OpenContractEvent(ctx cosmos.Context, contract types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeOpenContract,
				sdk.NewAttribute("pubkey", contract.ProviderPubKey.String()),
				sdk.NewAttribute("chain", contract.Chain.String()),
				sdk.NewAttribute("client", contract.Client.String()),
				sdk.NewAttribute("delegate", contract.Delegate.String()),
				sdk.NewAttribute("type", contract.Type.String()),
				sdk.NewAttribute("height", strconv.FormatInt(contract.Height, 10)),
				sdk.NewAttribute("duration", strconv.FormatInt(contract.Duration, 10)),
				sdk.NewAttribute("rate", strconv.FormatInt(contract.Rate, 10)),
			),
		},
	)
}

func (mgr Manager) ContractSettlementEvent(ctx cosmos.Context, debt cosmos.Int, contract types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeContractSettlement,
				sdk.NewAttribute("pubkey", contract.ProviderPubKey.String()),
				sdk.NewAttribute("chain", contract.Chain.String()),
				sdk.NewAttribute("client", contract.Client.String()),
				sdk.NewAttribute("paid", debt.String()),
			),
		},
	)
}

func (mgr Manager) ValidatorPayoutEvent(ctx cosmos.Context, acc cosmos.AccAddress, rwd cosmos.Int) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeValidatorPayout,
				sdk.NewAttribute("validator", acc.String()),
				sdk.NewAttribute("paid", rwd.String()),
			),
		},
	)
}
