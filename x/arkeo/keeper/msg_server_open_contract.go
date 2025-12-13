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
		"service", msg.Service,
		"client", msg.Client,
		"delegate", msg.Delegate,
		"contract type", msg.ContractType,
		"duration", msg.Duration,
		"rate", msg.Rate,
		"settlement duration", msg.SettlementDuration,
		"authorization", msg.Authorization,
	)

	// CacheContext implies NewEventManager
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

	service, _, err := k.ResolveServiceEnum(ctx, msg.Service)
	if err != nil {
		return err
	}
	providerPubKey, err := common.NewPubKey(msg.Provider)
	if err != nil {
		return err
	}
	provider, err := k.GetProvider(ctx, providerPubKey, service)
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

	switch msg.ContractType {
	case types.ContractType_SUBSCRIPTION:
		if cosmos.NewCoins(provider.SubscriptionRate...).AmountOf(msg.Rate.Denom).IsZero() {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "provider rates is 0, client sent %d", msg.Rate.Amount.Int64())
		}
		if !msg.Rate.Amount.Equal(cosmos.NewCoins(provider.SubscriptionRate...).AmountOf(msg.Rate.Denom)) {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "provider rates is %d, client sent %d", cosmos.NewCoins(provider.SubscriptionRate...).AmountOf(msg.Rate.Denom).Int64(), msg.Rate.Amount.Int64())
		}
		if !cosmos.NewInt(msg.Rate.Amount.Int64() * msg.Duration * msg.QueriesPerMinute).Equal(msg.Deposit) {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "mismatch of rate*duration*queriesPerMinute and deposit: %d * %d * %d != %d", msg.Rate.Amount.Int64(), msg.Duration, msg.QueriesPerMinute, msg.Deposit.Int64())
		}
	case types.ContractType_PAY_AS_YOU_GO:
		if cosmos.NewCoins(provider.PayAsYouGoRate...).AmountOf(msg.Rate.Denom).IsZero() {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "provider rates is 0, client sent %d", msg.Rate.Amount.Int64())
		}
		if !msg.Rate.Amount.Equal(cosmos.NewCoins(provider.PayAsYouGoRate...).AmountOf(msg.Rate.Denom)) {
			return errors.Wrapf(types.ErrOpenContractMismatchRate, "pay-as-you-go provider rate is %d, client sent %d", cosmos.NewCoins(provider.PayAsYouGoRate...).AmountOf(msg.Rate.Denom).Int64(), msg.Rate.Amount.Int64())
		}
		if msg.SettlementDuration != provider.SettlementDuration {
			return errors.Wrapf(types.ErrOpenContractMismatchSettlementDuration, "pay-as-you-go provider settlement duration is %d, client sent %d", provider.SettlementDuration, msg.SettlementDuration)
		}
	default:
		return errors.Wrapf(types.ErrInvalidContractType, "%s", msg.ContractType.String())
	}

	spender, err := msg.GetSpender()
	if err != nil {
		return err
	}

	activeContract, err := k.GetActiveContractForUser(ctx, spender, providerPubKey, service)
	if err != nil {
		return err
	}

	if !activeContract.IsEmpty() && !activeContract.IsExpired(ctx.BlockHeight()) {
		if msg.ContractType != types.ContractType_PAY_AS_YOU_GO {
			return errors.Wrapf(types.ErrOpenContractAlreadyOpen, "expires in %d blocks", activeContract.Expiration()-ctx.BlockHeight())
		}
		// For PAYG, allow multiple overlapping contracts
	}

	return nil
}

func (k msgServer) OpenContractHandle(ctx cosmos.Context, msg *types.MsgOpenContract) error {
	// set back client as delegate if delegate is empty
	if msg.Delegate == "" {
		msg.Delegate = msg.Client
	}

	openCost := k.FetchConfig(ctx, configs.OpenContractCost)
	if openCost > 0 {
		if err := k.SendFromAccountToModule(ctx, msg.MustGetSigner(), types.ReserveName, getCoins(openCost)); err != nil {
			return errors.Wrapf(err, "failed to send open contract costs openCost=%d", openCost)
		}
	}

	if err := k.SendFromAccountToModule(ctx, msg.MustGetSigner(), types.ContractName, cosmos.NewCoins(cosmos.NewCoin(msg.Rate.Denom, msg.Deposit))); err != nil {
		return errors.Wrapf(err, "failed to send deposit=%d", msg.Deposit.Int64())
	}

	service, svcRecord, err := k.ResolveServiceEnum(ctx, msg.Service)
	if err != nil {
		return err
	}

	providerPubKey, err := common.NewPubKey(msg.Provider)
	if err != nil {
		return types.ErrInvalidPubKey
	}

	clientPubKey, err := common.NewPubKey(msg.Client)
	if err != nil {
		return types.ErrInvalidPubKey
	}

	delegatePubKey, err := common.NewPubKey(msg.Delegate)
	if err != nil {
		return types.ErrInvalidPubKey
	}

	contract := types.Contract{
		Provider:           providerPubKey,
		Id:                 k.Keeper.GetAndIncrementNextContractId(ctx),
		Service:            service,
		Type:               msg.ContractType,
		Client:             clientPubKey,
		Delegate:           delegatePubKey,
		Duration:           msg.Duration,
		Rate:               msg.Rate,
		Deposit:            msg.Deposit,
		Paid:               cosmos.ZeroInt(),
		Height:             ctx.BlockHeight(),
		SettlementDuration: msg.SettlementDuration,
		Authorization:      msg.Authorization,
		QueriesPerMinute:   msg.QueriesPerMinute,
	}

	// create expiration set
	// these are used by the end blocker to settle contracts. We need to
	// use the additional settlement period for pay as you go contracts.
	expirationSet, err := k.GetContractExpirationSet(ctx, contract.SettlementPeriodEnd())
	if err != nil {
		return err
	}

	expirationSet.Append(contract.Id)
	err = k.SetContractExpirationSet(ctx, expirationSet)
	if err != nil {
		return err
	}

	spender, err := msg.GetSpender()
	if err != nil {
		return err
	}

	// create user set.
	userSet, err := k.GetUserContractSet(ctx, spender)
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

	ctx.Logger().Info("contract opened",
		"contract_id", contract.Id,
		"service", svcRecord.Name,
	)

	return k.EmitOpenContractEvent(ctx, openCost, &contract)
}
