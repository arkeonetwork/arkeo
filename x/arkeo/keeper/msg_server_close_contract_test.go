package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestCloseContractValidate(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(14)

	s := newMsgServer(k, sk)

	// setup
	providerPubkey := types.GetRandomPubKey()

	clientPubKey := types.GetRandomPubKey()
	clientAcct, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)
	service := common.BTCService

	contract := types.NewContract(providerPubkey, service, clientPubKey)
	contract.Duration = 100
	contract.Height = 10
	contract.Id = 1
	require.NoError(t, k.SetContract(ctx, contract))

	// happy path
	msg := types.MsgCloseContract{
		Creator:    clientAcct.String(),
		ContractId: contract.Id,
		Client:     clientPubKey,
	}
	require.NoError(t, s.CloseContractValidate(ctx, &msg))

	contract.Duration = 3
	require.NoError(t, k.SetContract(ctx, contract))
	err = s.CloseContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrCloseContractAlreadyClosed)
}

func TestCloseContractHandle(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	// setup
	providerPubKey := types.GetRandomPubKey()
	provider, err := providerPubKey.GetMyAddress()
	require.NoError(t, err)
	clientPubKey := types.GetRandomPubKey()
	clientAccount, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)

	service := common.BTCService
	require.True(t, k.GetBalance(ctx, provider).IsZero())
	rate, err := cosmos.ParseCoin("5uarkeo")
	require.NoError(t, err)

	openContractMessage := types.MsgOpenContract{
		Creator:      clientAccount.String(),
		Client:       clientPubKey.String(),
		Service:      service.String(),
		Provider:     providerPubKey.String(),
		Deposit:      cosmos.NewInt(500),
		Rate:         rate,
		Duration:     100,
		ContractType: types.ContractType_SUBSCRIPTION,
	}

	require.NoError(t, k.MintAndSendToAccount(ctx, clientAccount, getCoin(common.Tokens(10))))
	err = s.OpenContractHandle(ctx, &openContractMessage)
	require.NoError(t, err)

	bal := k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom)
	require.Equal(t, bal.Int64(), int64(500))

	contract, err := k.GetActiveContractForUser(ctx, clientPubKey, providerPubKey, service)
	require.NoError(t, err)
	require.False(t, contract.IsEmpty())

	ctx = ctx.WithBlockHeight(14)

	// happy path
	msg := types.MsgCloseContract{
		Creator:    clientAccount.String(),
		ContractId: contract.Id,
		Client:     clientPubKey,
	}
	require.NoError(t, s.CloseContractHandle(ctx, &msg))

	contract, err = k.GetContract(ctx, contract.Id)
	require.NoError(t, err)
	require.Equal(t, contract.Paid.Int64(), int64(20))
	require.Equal(t, contract.SettlementHeight, ctx.BlockHeight())

	bal = k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom)
	require.Equal(t, bal.Int64(), int64(0))
	require.True(t, k.HasCoins(ctx, provider, getCoins(18)))
	require.True(t, k.HasCoins(ctx, contract.ClientAddress(), getCoins(480)))
	bal = k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom)
	require.Equal(t, bal.Int64(), int64(100000002)) // open cost + fee
}

func TestCloseSubscriptionContract(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	creatroAddress := types.GetRandomBech32Addr()

	// set up provider
	providerPubKey := types.GetRandomPubKey()
	providerAddress, err := providerPubKey.GetMyAddress()
	require.NoError(t, err)
	service := common.BTCService
	provider := types.NewProvider(providerPubKey, service)
	provider.Bond = cosmos.NewInt(10000000000)
	require.NoError(t, k.SetProvider(ctx, provider))

	rates, err := cosmos.ParseCoins("15uarkeo")
	require.NoError(t, err)

	modProviderMsg := types.MsgModProvider{
		Creator:             creatroAddress.String(),
		Provider:            provider.PubKey,
		Service:             provider.Service.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      rates,
		SubscriptionRate:    rates,
	}
	err = s.ModProviderHandle(ctx, &modProviderMsg)

	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, providerAddress, getCoin(common.Tokens(10))))

	// set up contract with no delegate address
	userPubKey := types.GetRandomPubKey()
	userAddress, err := userPubKey.GetMyAddress()
	require.NoError(t, err)

	openContractMessage := types.MsgOpenContract{
		Provider:         providerPubKey.String(),
		Service:          service.String(),
		Creator:          providerAddress.String(),
		Client:           userPubKey.String(),
		ContractType:     types.ContractType_SUBSCRIPTION,
		Duration:         100,
		Rate:             rates[0],
		Deposit:          cosmos.NewInt(1500),
		QueriesPerMinute: 1,
	}
	_, err = s.OpenContract(ctx, &openContractMessage)
	require.NoError(t, err)

	contract, err := s.GetActiveContractForUser(ctx, userPubKey, providerPubKey, service)
	require.NoError(t, err)
	require.False(t, contract.IsEmpty())
	require.Equal(t, contract.Id, uint64(1))
	require.Equal(t, contract.Client, userPubKey)

	// confirm that another user cannot close the contract
	user2PubKey := types.GetRandomPubKey()
	user2Address, err := user2PubKey.GetMyAddress()
	require.NoError(t, err)

	closeContractMsg := types.MsgCloseContract{
		Creator:    user2Address.String(),
		ContractId: contract.Id,
		Client:     contract.Client,
	}
	_, err = s.CloseContract(ctx, &closeContractMsg)
	require.ErrorIs(t, err, types.ErrCloseContractUnauthorized)

	// confirm that the contract can be closed by the client
	closeContractMsg.Creator = userAddress.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	require.NoError(t, err)
	contract, err = s.GetActiveContractForUser(ctx, userPubKey, providerPubKey, service)
	require.NoError(t, err)
	require.True(t, contract.IsEmpty())

	// reopen contract this time with a delagate address.
	openContractMessage.Delegate = user2PubKey.String()
	_, err = s.OpenContract(ctx, &openContractMessage)
	require.NoError(t, err)

	contract, err = s.GetActiveContractForUser(ctx, user2PubKey, providerPubKey, service)
	require.NoError(t, err)
	require.Equal(t, contract.Id, uint64(2))
	require.False(t, contract.IsEmpty())
	require.True(t, contract.Delegate.Equals(user2PubKey))
	closeContractMsg.ContractId = contract.Id

	user2ContractSet, err := s.GetUserContractSet(ctx, user2PubKey)
	require.NoError(t, err)
	require.Len(t, user2ContractSet.ContractSet.ContractIds, 1)

	// confirm that the contract cannot be closed by the delegate
	closeContractMsg.Creator = user2Address.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	require.ErrorIs(t, err, types.ErrCloseContractUnauthorized)

	// but can be closed by the client
	closeContractMsg.Creator = userAddress.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	require.NoError(t, err)
}

func TestClosePayAsYouGoContract(t *testing.T) {
	// NOTE: pay as you go contracts cannot be closed on demand.
	ctx, k, sk := SetupKeeperWithStaking(t)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	// set up provider
	providerPubKey := types.GetRandomPubKey()
	providerAddress, err := providerPubKey.GetMyAddress()
	require.NoError(t, err)
	service := common.BTCService
	provider := types.NewProvider(providerPubKey, service)
	provider.Bond = cosmos.NewInt(10000000000)
	require.NoError(t, k.SetProvider(ctx, provider))

	rates, err := cosmos.ParseCoins("15uarkeo")
	require.NoError(t, err)

	creatorAddress := types.GetRandomBech32Addr()

	modProviderMsg := types.MsgModProvider{
		Creator:             creatorAddress.String(),
		Provider:            provider.PubKey,
		Service:             provider.Service.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      rates,
		SubscriptionRate:    rates,
	}
	err = s.ModProviderHandle(ctx, &modProviderMsg)

	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, providerAddress, getCoin(common.Tokens(10))))

	// set up contract with no delegate address
	userPubKey := types.GetRandomPubKey()
	userAddress, err := userPubKey.GetMyAddress()
	require.NoError(t, err)

	openContractMessage := types.MsgOpenContract{
		Provider:     providerPubKey.String(),
		Service:      service.String(),
		Creator:      providerAddress.String(),
		Client:       userPubKey.String(),
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         rates[0],
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &openContractMessage)
	require.NoError(t, err)

	contract, err := s.GetActiveContractForUser(ctx, userPubKey, providerPubKey, service)
	require.NoError(t, err)

	// confirm that another user cannot close the contract
	use2PubKey := types.GetRandomPubKey()
	user2Address, err := use2PubKey.GetMyAddress()
	require.NoError(t, err)

	closeContractMsg := types.MsgCloseContract{
		Creator:    user2Address.String(),
		ContractId: contract.Id,
		Client:     contract.Client,
	}
	_, err = s.CloseContract(ctx, &closeContractMsg)
	require.ErrorIs(t, err, types.ErrCloseContractUnauthorized)

	// reopen contract this time with a delagate address.
	openContractMessage.Delegate = use2PubKey.String()
	_, err = s.OpenContract(ctx, &openContractMessage)
	require.NoError(t, err)

	contract, err = s.GetActiveContractForUser(ctx, userPubKey, providerPubKey, service)
	require.NoError(t, err)

	closeContractMsg.ContractId = contract.Id

	// confirm that the contract cannot be closed by the delegate
	closeContractMsg.Creator = user2Address.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	require.ErrorIs(t, err, types.ErrCloseContractUnauthorized)

	// unbond provider , unbond 100% will remove the provider
	closeContractMsg.Creator = userAddress.String()
	k.RemoveProvider(ctx, providerPubKey, service)
	_, err = s.CloseContract(ctx, &closeContractMsg)
	require.NoError(t, err)
}

func TestCloseContractUnauthorizedClient(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)

	ctx = ctx.WithBlockHeight(14)

	s := newMsgServer(k, sk)

	// setup

	providerPubKey := types.GetRandomPubKey()

	clientPubKey := types.GetRandomPubKey()

	_, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)
	serivce := common.BTCService

	contract := types.NewContract(providerPubKey, serivce, clientPubKey)

	contract.Duration = 100
	contract.Height = 10
	contract.Id = 1
	require.NoError(t, k.SetContract(ctx, contract))

	unauthorizedClientPubKey := types.GetRandomPubKey()
	unauthorizedClientAddress, err := unauthorizedClientPubKey.GetMyAddress()

	require.NoError(t, err)
	msg := types.MsgCloseContract{
		Creator:    unauthorizedClientAddress.String(),
		ContractId: contract.Id,
		Client:     unauthorizedClientPubKey,
	}

	err = s.CloseContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrCloseContractUnauthorized)
}

func TestCloseContractWithIncorrectDelegate(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)

	ctx = ctx.WithBlockHeight(14)

	s := newMsgServer(k, sk)

	// setup

	providerPubKey := types.GetRandomPubKey()

	clientPubKey := types.GetRandomPubKey()

	clientAcct, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)
	serivce := common.BTCService

	contract := types.NewContract(providerPubKey, serivce, clientPubKey)

	contract.Duration = 100
	contract.Height = 10
	contract.Id = 1
	contract.Delegate = types.GetRandomPubKey()
	require.NoError(t, k.SetContract(ctx, contract))

	require.NoError(t, err)
	msg := types.MsgCloseContract{
		Creator:    clientAcct.String(),
		ContractId: contract.Id,
		Client:     clientPubKey,
		Delegate:   types.GetRandomPubKey(),
	}

	err = s.CloseContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrCloseContractUnauthorized)
}

func TestCloseContractWithNoClient(t *testing.T) {
	ctx, k, sk := SetupKeeperWithStaking(t)

	ctx = ctx.WithBlockHeight(14)

	s := newMsgServer(k, sk)

	// setup

	providerPubKey := types.GetRandomPubKey()
	clientPubKey := types.GetRandomPubKey()

	_, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)
	serivce := common.BTCService

	contract := types.NewContract(providerPubKey, serivce, clientPubKey)
	contract.Duration = 100
	contract.Height = 10
	contract.Id = 1
	require.NoError(t, k.SetContract(ctx, contract))

	unauthorizedClientPubKey := types.GetRandomPubKey()
	unauthorizedClientAddress, err := unauthorizedClientPubKey.GetMyAddress()

	require.NoError(t, err)
	msg := types.MsgCloseContract{
		Creator:    unauthorizedClientAddress.String(),
		ContractId: contract.Id,
	}

	err = s.CloseContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrCloseContractUnauthorized)
}
