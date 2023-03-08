package keeper

import (
	"strconv"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) BondProviderEvent(ctx cosmos.Context, bond cosmos.Int, msg *types.MsgBondProvider) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeProviderBond,
				sdk.NewAttribute("provider", msg.Provider.String()),
				sdk.NewAttribute("chain", msg.Chain),
				sdk.NewAttribute("bond_rel", msg.Bond.String()),
				sdk.NewAttribute("bond_abs", bond.String()),
			),
		},
	)
}

func (k msgServer) CloseContractEvent(ctx cosmos.Context, contract *types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeCloseContract,
				sdk.NewAttribute("contract_id", strconv.FormatUint(contract.Id, 10)),
				sdk.NewAttribute("provider", contract.Provider.String()),
				sdk.NewAttribute("chain", contract.Chain.String()),
				sdk.NewAttribute("client", contract.Client.String()),
				sdk.NewAttribute("delegate", contract.Delegate.String()),
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
				sdk.NewAttribute("metadata_uri", provider.MetadataUri),
				sdk.NewAttribute("metadata_nonce", strconv.FormatUint(provider.MetadataNonce, 10)),
				sdk.NewAttribute("status", provider.Status.String()),
				sdk.NewAttribute("min_contract_duration", strconv.FormatInt(provider.MinContractDuration, 10)),
				sdk.NewAttribute("max_contract_duration", strconv.FormatInt(provider.MaxContractDuration, 10)),
				sdk.NewAttribute("subscription_rate", strconv.FormatInt(provider.SubscriptionRate, 10)),
				sdk.NewAttribute("pay-as-you-go_rate", strconv.FormatInt(provider.PayAsYouGoRate, 10)),
				sdk.NewAttribute("bond", provider.Bond.String()),
				sdk.NewAttribute("settlement_duration", strconv.FormatInt(provider.SettlementDuration, 10)),
			),
		},
	)
}

func (k msgServer) OpenContractEvent(ctx cosmos.Context, openCost int64, contract types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeOpenContract,
				sdk.NewAttribute("provider", contract.Provider.String()),
				sdk.NewAttribute("contract_id", strconv.FormatUint(contract.Id, 10)),
				sdk.NewAttribute("chain", contract.Chain.String()),
				sdk.NewAttribute("client", contract.Client.String()),
				sdk.NewAttribute("delegate", contract.Delegate.String()),
				sdk.NewAttribute("type", contract.Type.String()),
				sdk.NewAttribute("height", strconv.FormatInt(contract.Height, 10)),
				sdk.NewAttribute("duration", strconv.FormatInt(contract.Duration, 10)),
				sdk.NewAttribute("rate", strconv.FormatInt(contract.Rate, 10)),
				sdk.NewAttribute("open_cost", strconv.FormatInt(openCost, 10)),
				sdk.NewAttribute("deposit", contract.Deposit.String()),
				sdk.NewAttribute("settlement_duration", strconv.FormatInt(contract.SettlementDuration, 10)),
			),
		},
	)
}

func (mgr Manager) ContractSettlementEvent(ctx cosmos.Context, debt, valIncome cosmos.Int, contract types.Contract) {
	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeContractSettlement,
				sdk.NewAttribute("provider", contract.Provider.String()),
				sdk.NewAttribute("contract_id", strconv.FormatUint(contract.Id, 10)),
				sdk.NewAttribute("chain", contract.Chain.String()),
				sdk.NewAttribute("client", contract.Client.String()),
				sdk.NewAttribute("delegate", contract.Delegate.String()),
				sdk.NewAttribute("type", contract.Type.String()),
				sdk.NewAttribute("nonce", strconv.FormatInt(contract.Nonce, 10)),
				sdk.NewAttribute("height", strconv.FormatInt(contract.Height, 10)),
				sdk.NewAttribute("paid", debt.String()),
				sdk.NewAttribute("reserve", valIncome.String()),
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
