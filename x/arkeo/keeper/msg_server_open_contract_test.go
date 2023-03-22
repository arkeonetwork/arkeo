package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestOpenContractValidate(t *testing.T) {
	var err error
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(1)
	s := newMsgServer(k, sk)

	// setup
	providerPubkey := types.GetRandomPubKey()
	clientPubKey := types.GetRandomPubKey()
	acc, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)

	service := common.BTCService

	provider := types.NewProvider(providerPubkey, service)
	provider.Bond = cosmos.NewInt(500_00000000)
	provider.Status = types.ProviderStatus_ONLINE
	provider.MaxContractDuration = 1000
	provider.MinContractDuration = 10
	provider.SubscriptionRate = 15
	provider.PayAsYouGoRate = 2
	provider.LastUpdate = 1
	provider.PayAsYouGoEnabled = true
	provider.SubscriptionEnabled = true
	require.NoError(t, k.SetProvider(ctx, provider))

	// happy path
	msg := types.MsgOpenContract{
		Provider:     providerPubkey,
		Service:      service.String(),
		Client:       clientPubKey,
		Creator:      acc.String(),
		ContractType: types.ContractType_SUBSCRIPTION,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(100 * 15),
	}
	require.NoError(t, k.MintAndSendToAccount(ctx, acc, getCoin(common.Tokens(100*25))))
	require.NoError(t, s.OpenContractValidate(ctx, &msg))

	// check duration
	msg.Duration = 10000000000000
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractDuration)
	msg.Duration = 5
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractDuration)
	msg.Duration = 100

	// check rates
	msg.Rate = 10
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractMismatchRate)
	msg.ContractType = types.ContractType_PAY_AS_YOU_GO
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractMismatchRate)
	msg.Rate = 15
	msg.ContractType = types.ContractType_SUBSCRIPTION

	provider.Bond = cosmos.NewInt(1)
	require.NoError(t, k.SetProvider(ctx, provider))
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrInvalidBond)
	provider.Bond = cosmos.NewInt(500_00000000)
	require.NoError(t, k.SetProvider(ctx, provider))

	ctx = ctx.WithBlockHeight(15)
	contract := types.NewContract(providerPubkey, service, clientPubKey)
	contract.Type = types.ContractType_SUBSCRIPTION
	contract.Height = ctx.BlockHeight()
	contract.Duration = 100
	contract.Rate = 2

	require.NoError(t, s.OpenContractHandle(ctx, &msg))
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractAlreadyOpen)
}

func TestOpenContractHandle(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	// setup
	pubkey := types.GetRandomPubKey()
	acc, err := pubkey.GetMyAddress()
	require.NoError(t, err)
	service := common.BTCService
	require.NoError(t, k.MintAndSendToAccount(ctx, acc, getCoin(common.Tokens(10))))

	msg := types.MsgOpenContract{
		Provider:     pubkey,
		Service:      service.String(),
		Creator:      acc.String(),
		Client:       pubkey,
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1000),
	}
	require.NoError(t, s.OpenContractHandle(ctx, &msg))

	contract, err := k.GetActiveContractForUser(ctx, pubkey, pubkey, service)
	require.NoError(t, err)

	require.Equal(t, contract.Type, types.ContractType_PAY_AS_YOU_GO)
	require.False(t, contract.IsEmpty())

	require.Equal(t, contract.Height, ctx.BlockHeight())
	require.Equal(t, contract.Duration, int64(100))
	require.Equal(t, contract.Rate, int64(15))
	require.Equal(t, contract.Nonce, int64(0))
	require.Equal(t, contract.Deposit.Int64(), int64(1000))
	require.Equal(t, contract.Paid.Int64(), int64(0))

	bal := k.GetBalance(ctx, acc) // check balance
	require.Equal(t, bal.AmountOf(configs.Denom).Int64(), int64(899999000))

	// check that contract expiration has been set
	set, err := k.GetContractExpirationSet(ctx, contract.Expiration())
	require.NoError(t, err)
	require.Equal(t, set.Height, contract.Expiration())
	require.Len(t, set.ContractSet.ContractIds, 1)

	// check that contract has been added to the user
	userSet, err := k.GetUserContractSet(ctx, contract.GetSpender())
	require.NoError(t, err)
	require.Equal(t, userSet.User, contract.GetSpender())
	require.Len(t, userSet.ContractSet.ContractIds, 1)
}

func TestOpenContract(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	providerPubKey := types.GetRandomPubKey()
	providerAddress, err := providerPubKey.GetMyAddress()
	require.NoError(t, err)
	service := common.BTCService
	provider := types.NewProvider(providerPubKey, service)
	provider.Bond = cosmos.NewInt(10000000000)
	provider.LastUpdate = ctx.BlockHeight()
	require.NoError(t, k.SetProvider(ctx, provider))

	modProviderMsg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Service:             provider.Service.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      15,
		SubscriptionRate:    15,
		PayAsYouGoEnabled:   true,
		SubscriptionEnabled: true,
	}
	err = s.ModProviderHandle(ctx, &modProviderMsg)

	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, providerAddress, getCoin(common.Tokens(10))))

	msg := types.MsgOpenContract{
		Provider:     providerPubKey,
		Service:      service.String(),
		Creator:      providerAddress.String(),
		Client:       providerPubKey,
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	contract, err := k.GetActiveContractForUser(ctx, providerPubKey, providerPubKey, service)
	require.NoError(t, err)

	require.False(t, contract.IsEmpty())
	require.Equal(t, contract.Id, uint64(0))

	clientPubKey := types.GetRandomPubKey()
	clientAddress, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)

	msg = types.MsgOpenContract{
		Provider:     providerPubKey,
		Service:      service.String(),
		Creator:      clientAddress.String(),
		Client:       clientPubKey,
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1000),
	}
	require.NoError(t, k.MintAndSendToAccount(ctx, clientAddress, getCoin(common.Tokens(10))))
	require.NoError(t, s.OpenContractHandle(ctx, &msg))

	contract, err = k.GetActiveContractForUser(ctx, clientPubKey, providerPubKey, service)
	require.NoError(t, err)

	require.False(t, contract.IsEmpty())
	require.Equal(t, contract.Id, uint64(1))

	_, err = s.OpenContract(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractAlreadyOpen)

	// confirm that the client can open a contract with a deleagate
	delegatePubKey := types.GetRandomPubKey()
	msg.Delegate = delegatePubKey
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	contract, err = k.GetActiveContractForUser(ctx, delegatePubKey, providerPubKey, service)
	require.NoError(t, err)

	require.False(t, contract.IsEmpty())
	require.Equal(t, contract.Id, uint64(2))
}

func TestOpenContractWithSettlementPeriod(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	providerPubKey := types.GetRandomPubKey()
	service := common.BTCService
	provider := types.NewProvider(providerPubKey, service)
	provider.Bond = cosmos.NewInt(10000000000)
	provider.LastUpdate = ctx.BlockHeight()
	require.NoError(t, k.SetProvider(ctx, provider))

	modProviderMsg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Service:             provider.Service.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      15,
		SubscriptionRate:    15,
		SettlementDuration:  10,
		PayAsYouGoEnabled:   true,
		SubscriptionEnabled: true,
	}
	err := s.ModProviderHandle(ctx, &modProviderMsg)

	require.NoError(t, err)

	clientPubKey := types.GetRandomPubKey()
	clientAddress, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, clientAddress, getCoin(common.Tokens(10))))

	msg := types.MsgOpenContract{
		Provider:     providerPubKey,
		Service:      service.String(),
		Creator:      clientAddress.String(),
		Client:       clientPubKey,
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractMismatchSettlementDuration)

	msg.SettlementDuration = 10
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	contract, err := k.GetActiveContractForUser(ctx, clientPubKey, providerPubKey, service)
	require.NoError(t, err)

	require.False(t, contract.IsEmpty())
	require.Equal(t, contract.Id, uint64(0))

	// confirm opening a new contract will fail since the user have an active one
	_, err = s.OpenContract(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractAlreadyOpen)

	// move to a block where out contract should be expired, but not settled.
	ctx = ctx.WithBlockHeight(contract.Expiration() + 1)

	require.True(t, contract.IsExpired(ctx.BlockHeight()))
	require.True(t, contract.IsSettlementPeriod(ctx.BlockHeight()))
	require.False(t, contract.IsSettled(ctx.BlockHeight()))

	// client should be able to open another contract.
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	// confirm contract income can be claimed while the first contract is in the
	// settlement period.
	claimMsg := types.MsgClaimContractIncome{
		ContractId: contract.Id,
		Creator:    clientAddress.String(),
		Spender:    clientPubKey,
		Nonce:      20,
	}
	_, err = s.ClaimContractIncome(ctx, &claimMsg)
	require.NoError(t, err)

	// advance beyond settlement period to confirm the contract is settled and no more
	// income can be claimed.
	ctx = ctx.WithBlockHeight(contract.SettlementPeriodEnd())
	require.True(t, contract.IsExpired(ctx.BlockHeight()))
	require.False(t, contract.IsSettlementPeriod(ctx.BlockHeight()))
	require.True(t, contract.IsSettled(ctx.BlockHeight()))

	claimMsg.Nonce = 21
	_, err = s.ClaimContractIncome(ctx, &claimMsg)
	require.ErrorIs(t, err, types.ErrClaimContractIncomeClosed)
}

func TestOpenContractNotSupportPayAsYouGo(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	providerPubKey := types.GetRandomPubKey()
	service := common.BTCService
	provider := types.NewProvider(providerPubKey, service)
	provider.Bond = cosmos.NewInt(10000000000)
	provider.LastUpdate = ctx.BlockHeight()
	require.NoError(t, k.SetProvider(ctx, provider))

	modProviderMsg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Service:             provider.Service.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      15,
		SubscriptionRate:    15,
		SettlementDuration:  10,
		PayAsYouGoEnabled:   false,
		SubscriptionEnabled: true,
	}
	err := s.ModProviderHandle(ctx, &modProviderMsg)
	require.NoError(t, err)

	clientPubKey := types.GetRandomPubKey()
	clientAddress, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, clientAddress, getCoin(common.Tokens(10))))

	msg := types.MsgOpenContract{
		Provider:           providerPubKey,
		Service:            service.String(),
		Creator:            clientAddress.String(),
		Client:             clientPubKey,
		ContractType:       types.ContractType_PAY_AS_YOU_GO,
		Duration:           100,
		Rate:               15,
		Deposit:            cosmos.NewInt(1500),
		SettlementDuration: 10,
	}
	_, err = s.OpenContract(ctx, &msg)
	require.ErrorIs(t, err, types.ErrProviderNotSupportPayAsYouGo)
}
