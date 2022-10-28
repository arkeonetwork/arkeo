package keeper

import (
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Manager struct {
	keeper  Keeper
	configs configs.ConfigValues
}

func NewManager(k Keeper) Manager {
	ver := k.GetVersion()
	return Manager{
		keeper:  k,
		configs: configs.GetConfigValues(ver),
	}
}

func (mgr Manager) EndBlock(ctx cosmos.Context) error {
	if err := mgr.ContractEndBlock(ctx); err != nil {
		ctx.Logger().Error("unable to settle contracts", "error", err)
	}
	return nil
}

func (mgr Manager) ContractEndBlock(ctx cosmos.Context) error {
	set, err := mgr.keeper.GetContractExpirationSet(ctx, ctx.BlockHeight())
	if err != nil {
		return err
	}

	for _, exp := range set.Contracts {
		contract, err := mgr.keeper.GetContract(ctx, exp.ProviderPubKey, exp.Chain, exp.ClientAddress)
		if err != nil {
			ctx.Logger().Error("unable to fetch contract", "pubkey", exp.ProviderPubKey, "chain", exp.Chain, "client", exp.ClientAddress, "error", err)
			continue
		}
		_, err = mgr.SettleContract(ctx, contract, true)
		if err != nil {
			ctx.Logger().Error("unable settle contract", "pubkey", exp.ProviderPubKey, "chain", exp.Chain, "client", exp.ClientAddress, "error", err)
			continue
		}
	}

	return nil
}

// any owed debt is paid to data provider
func (mgr Manager) SettleContract(ctx cosmos.Context, contract types.Contract, closed bool) (types.Contract, error) {
	debt, err := mgr.contractDebt(ctx, contract)
	if err != nil {
		return contract, err
	}
	if !debt.IsZero() {
		provider, err := contract.ProviderPubKey.GetMyAddress()
		if err != nil {
			return contract, err
		}
		if err := mgr.keeper.SendFromModuleToAccount(ctx, types.ContractName, provider, cosmos.NewCoins(cosmos.NewCoin(configs.Denom, debt))); err != nil {
			return contract, err
		}
	}

	contract.Paid = contract.Paid.Add(debt)
	if closed {
		remainder := contract.Deposit.Sub(contract.Paid)
		if !remainder.IsZero() {
			if err := mgr.keeper.SendFromModuleToAccount(ctx, types.ContractName, contract.ClientAddress, cosmos.NewCoins(cosmos.NewCoin(configs.Denom, remainder))); err != nil {
				return contract, err
			}
		}
		contract.ClosedHeight = ctx.BlockHeight()
	}

	err = mgr.keeper.SetContract(ctx, contract)
	if err != nil {
		return contract, err
	}

	ctx.EventManager().EmitEvents(
		sdk.Events{
			sdk.NewEvent(
				types.EventTypeContractSettlement,
				sdk.NewAttribute("pubkey", contract.ProviderPubKey.String()),
				sdk.NewAttribute("chain", contract.Chain.String()),
				sdk.NewAttribute("client", contract.ClientAddress.String()),
				sdk.NewAttribute("paid", debt.String()),
			),
		},
	)
	return contract, nil
}

func (mgr Manager) contractDebt(ctx cosmos.Context, contract types.Contract) (cosmos.Int, error) {
	var debt cosmos.Int
	switch contract.Type {
	case types.ContractType_Subscription:
		debt = cosmos.NewInt(contract.Rate * (ctx.BlockHeight() - contract.Height)).Sub(contract.Paid)
	case types.ContractType_PayAsYouGo:
		debt = cosmos.NewInt(contract.Rate * contract.Queries).Sub(contract.Paid)
	default:
		return cosmos.ZeroInt(), sdkerrors.Wrapf(types.ErrInvalidContractType, "%s", contract.Type.String())
	}

	if debt.IsNegative() {
		return cosmos.ZeroInt(), nil
	}
	return debt, nil
}
