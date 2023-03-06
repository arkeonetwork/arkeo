package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (k msgServer) OpenContract(goCtx context.Context, msg *types.MsgOpenContract) (*types.MsgOpenContractResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgOpenContract",
		"provder", msg.Provider,
		"chain", msg.Chain,
		"client", msg.Client,
		"delegate", msg.Delegate,
		"contract type", msg.ContractType,
		"duration", msg.Duration,
		"rate", msg.Rate,
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.OpenContractValidate(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed open contract validation", "err", err)
		return nil, err
	}

	if err := k.OpenContractHandle(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed open contract handle", "err", err)
		return nil, err
	}
	commit()

	return &types.MsgOpenContractResponse{}, nil
}

func (k msgServer) OpenContractValidate(ctx cosmos.Context, msg *types.MsgOpenContract) error {
	if k.FetchConfig(ctx, configs.HandlerOpenContract) > 0 {
		return errors.Wrapf(types.ErrDisabledHandler, "open contract")
	}

	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}
	provider, err := k.GetProvider(ctx, msg.Provider, chain)
	if err != nil {
		return err
	}

	if provider.LastUpdate == 0 {
		return errors.Wrapf(types.ErrProviderNotFound, "provider %s for chain %s not found", msg.Provider, msg.Chain)
	}

	minBond := k.FetchConfig(ctx, configs.MinProviderBond)
	if provider.Bond.LT(cosmos.NewInt(minBond)) {
		return errors.Wrapf(types.ErrInvalidBond, "not enough provider bond to open a contract (%d/%d)", provider.Bond.Int64(), minBond)
	}

	if provider.Status != types.ProviderStatus_ONLINE {
		return errors.Wrapf(types.ErrOpenContractBadProviderStatus, "has status %s", provider.Status.String())
	}

	if msg.Duration > provider.MaxContractDuration {
		return errors.Wrapf(types.ErrOpenContractDuration, "duration exceeds allowed maximum duration from provider")
	}

	if msg.Duration < provider.MinContractDuration {
		return errors.Wrapf(types.ErrOpenContractDuration, "duration below allowed minimum duration from provider")
	}

	switch msg.ContractType {
	case types.ContractType_SUBSCRIPTION:
		if msg.Rate != provider.SubscriptionRate {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "provider rates is %d, client sent %d", provider.SubscriptionRate, msg.Rate)
		}
		if !cosmos.NewInt(msg.Rate * msg.Duration).Equal(msg.Deposit) {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "mismatch of rate*duration and deposit: %d * %d != %d", msg.Rate, msg.Duration, msg.Deposit.Int64())
		}
	case types.ContractType_PAY_AS_YOU_GO:
		if msg.Rate != provider.PayAsYouGoRate {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "pay-as-you-go provider rate is %d, client sent %d", provider.PayAsYouGoRate, msg.Rate)
		}
	default:
		return errors.Wrapf(types.ErrInvalidContractType, "%s", msg.ContractType.String())
	}

	activeContract, err := k.GetActiveContractForUser(ctx, msg.GetSpender(), msg.Provider, chain)
	if err != nil {
		return err
	}

	if !activeContract.IsEmpty() {
		return errors.Wrapf(types.ErrOpenContractAlreadyOpen, "expires in %d blocks", activeContract.Expiration()-ctx.BlockHeight())
	}

	return nil
}

func (k msgServer) OpenContractHandle(ctx cosmos.Context, msg *types.MsgOpenContract) error {
	openCost := k.FetchConfig(ctx, configs.OpenContractCost)
	if openCost > 0 {
		if err := k.SendFromAccountToModule(ctx, msg.MustGetSigner(), types.ReserveName, getCoins(openCost)); err != nil {
			return errors.Wrapf(err, "failed to send open contract costs openCost=%d", openCost)
		}
	}

	if err := k.SendFromAccountToModule(ctx, msg.MustGetSigner(), types.ContractName, getCoins(msg.Deposit.Int64())); err != nil {
		return errors.Wrapf(err, "failed to send deposit=%d", msg.Deposit.Int64())
	}

	chain, err := common.NewChain(msg.Chain)
	if err != nil {
		return err
	}

	contract := types.NewContract(msg.Provider, chain, msg.GetSpender())
	contract.Id = k.Keeper.GetAndIncrementNextContractId(ctx)
	contract.Client = msg.Client
	contract.Type = msg.ContractType
	contract.Height = ctx.BlockHeight()
	contract.Duration = msg.Duration
	contract.Rate = msg.Rate
	contract.Deposit = msg.Deposit

	// create expiration set

	expirationSet, err := k.GetContractExpirationSet(ctx, contract.Expiration())
	if err != nil {
		return err
	}

	if expirationSet.ContractSet == nil {
		expirationSet.ContractSet = &types.ContractSet{}
	}

	expirationSet.ContractSet.ContractIds = append(expirationSet.ContractSet.ContractIds, contract.Id)
	err = k.SetContractExpirationSet(ctx, expirationSet)
	if err != nil {
		return err
	}

	// create user set.
	userSet, err := k.GetUserContractSet(ctx, msg.GetSpender())
	if err != nil {
		return err
	}

	if userSet.ContractSet == nil {
		userSet.ContractSet = &types.ContractSet{}
	}

	userSet.ContractSet.ContractIds = append(userSet.ContractSet.ContractIds, contract.Id)
	err = k.SetUserContractSet(ctx, userSet)
	if err != nil {
		return err
	}

	err = k.SetContract(ctx, contract)
	if err != nil {
		return err
	}

	k.OpenContractEvent(ctx, openCost, contract)
	return nil
}
