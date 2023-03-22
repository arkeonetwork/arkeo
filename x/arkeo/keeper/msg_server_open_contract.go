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
		"provider", msg.Provider,
		"service", msg.Service,
		"client", msg.Client,
		"delegate", msg.Delegate,
		"user type", msg.UserType,
		"meter type", msg.MeterType,
		"duration", msg.Duration,
		"rate", msg.Rate,
		"settlement duration", msg.SettlementDuration,
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

	service, err := common.NewService(msg.Service)
	if err != nil {
		return err
	}
	provider, err := k.GetProvider(ctx, msg.Provider, service)
	if err != nil {
		return err
	}

	if provider.LastUpdate == 0 {
		return errors.Wrapf(types.ErrProviderNotFound, "provider %s for service %s not found", msg.Provider, msg.Service)
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

	providerRate := types.FindRate(provider.Rates, msg.UserType, msg.MeterType)
	if providerRate == nil {
		return errors.Wrapf(types.ErrOpenContractMismatchRate, "provider does not currently support MeterType:%s for UserType:%s", msg.MeterType, msg.UserType)
	}

	switch msg.MeterType {
	case types.MeterType_PAY_PER_BLOCK:
		if msg.Rate != providerRate.Rate {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "provider rates is %d, client sent %d", providerRate.Rate, msg.Rate)
		}
		rate := cosmos.NewInt(msg.Rate)
		duration := cosmos.NewInt(msg.Duration) // we have confirmed that duration is positive in validate basic
		durationRate := rate.Mul(duration)
		if !durationRate.Equal(msg.Deposit) {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "mismatch of rate*duration and deposit: %d * %d != %d", msg.Rate, msg.Duration, msg.Deposit.Uint64())
		}
	case types.MeterType_PAY_PER_CALL:
		if msg.Rate != providerRate.Rate {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "provider rate for MeterType:%s for UserType:%s is %d, client sent %d", msg.MeterType, msg.UserType, providerRate.Rate, msg.Rate)
		}
		if msg.SettlementDuration != provider.SettlementDuration {
			return errors.Wrapf(types.ErrOpenContractMismatchSettlementDuration, "Pay-per-call provider settlement duration is %d, client sent %d", provider.SettlementDuration, msg.SettlementDuration)
		}
	default:
		return errors.Wrapf(types.ErrInvalidMeterType, "%s", msg.MeterType.String())
	}

	activeContract, err := k.GetActiveContractForUser(ctx, msg.GetSpender(), msg.Provider, service)
	if err != nil {
		return err
	}

	if !activeContract.IsEmpty() && !activeContract.IsExpired(ctx.BlockHeight()) {
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
		return errors.Wrapf(err, "failed to send deposit=%d", msg.Deposit.Uint64())
	}

	service, err := common.NewService(msg.Service)
	if err != nil {
		return err
	}

	contract := types.Contract{
		Provider:           msg.Provider,
		Id:                 k.Keeper.GetAndIncrementNextContractId(ctx),
		Service:            service,
		UserType:           msg.UserType,
		MeterType:          msg.MeterType,
		Client:             msg.Client,
		Delegate:           msg.Delegate,
		Duration:           msg.Duration,
		Rate:               msg.Rate,
		Deposit:            msg.Deposit,
		Paid:               cosmos.ZeroInt(),
		Height:             ctx.BlockHeight(),
		SettlementDuration: msg.SettlementDuration,
	}

	// create expiration set
	// these are used by the end blocker to settle contracts. We need to
	// use the additional settlement period for pay as you go contracts.
	expirationSet, err := k.GetContractExpirationSet(ctx, contract.SettlementPeriodEnd())
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

	k.OpenContractEvent(ctx, openCost, &contract)
	return nil
}
