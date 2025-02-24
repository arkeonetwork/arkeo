package keeper

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	cKeys "github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/std"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
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

	sRates, err := cosmos.ParseCoins("15uarkeo,20uatom")
	require.NoError(t, err)
	pRates, err := cosmos.ParseCoins("2uarkeo")
	require.NoError(t, err)

	provider := types.NewProvider(providerPubkey, service)
	provider.Bond = cosmos.NewInt(500_00000000)
	provider.Status = types.ProviderStatus_ONLINE
	provider.MaxContractDuration = 1000
	provider.MinContractDuration = 10
	provider.SubscriptionRate = sRates
	provider.PayAsYouGoRate = pRates
	provider.LastUpdate = 1
	require.NoError(t, k.SetProvider(ctx, provider))

	// happy path
	msg := types.MsgOpenContract{
		Provider:         providerPubkey.String(),
		Service:          service.String(),
		Client:           clientPubKey.String(),
		Creator:          acc.String(),
		ContractType:     types.ContractType_SUBSCRIPTION,
		Duration:         100,
		Rate:             sRates[0],
		Deposit:          cosmos.NewInt(100 * 15),
		QueriesPerMinute: 1,
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
	msg.Rate = cosmos.NewInt64Coin("bogus", 10)
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractMismatchRate)
	msg.Rate = cosmos.NewInt64Coin("uarkeo", 10)
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractMismatchRate)
	msg.Rate = cosmos.NewInt64Coin("uatom", 10)
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractMismatchRate)
	msg.ContractType = types.ContractType_PAY_AS_YOU_GO
	err = s.OpenContractValidate(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractMismatchRate)
	msg.Rate = cosmos.NewInt64Coin("uarkeo", 15)
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
	contract.Rate = pRates[0]

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
		Provider:         pubkey.String(),
		Service:          service.String(),
		Creator:          acc.String(),
		Client:           pubkey.String(),
		ContractType:     types.ContractType_PAY_AS_YOU_GO,
		Duration:         100,
		Rate:             cosmos.NewInt64Coin("uarkeo", 15),
		Deposit:          cosmos.NewInt(1000),
		QueriesPerMinute: 1,
	}
	require.NoError(t, s.OpenContractHandle(ctx, &msg))

	contract, err := k.GetActiveContractForUser(ctx, pubkey, pubkey, service)
	require.NoError(t, err)

	require.Equal(t, contract.Type, types.ContractType_PAY_AS_YOU_GO)
	require.False(t, contract.IsEmpty())

	require.Equal(t, contract.Height, ctx.BlockHeight())
	require.Equal(t, contract.Duration, int64(100))
	require.Equal(t, contract.Rate.Amount.Int64(), int64(15))
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

	msg := types.MsgOpenContract{
		Provider:     providerPubKey.String(),
		Service:      service.String(),
		Creator:      providerAddress.String(),
		Client:       providerPubKey.String(),
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         cosmos.NewInt64Coin("uarkeo", 15),
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	contract, err := k.GetActiveContractForUser(ctx, providerPubKey, providerPubKey, service)
	require.NoError(t, err)

	require.False(t, contract.IsEmpty())
	require.Equal(t, contract.Id, uint64(1))

	clientPubKey := types.GetRandomPubKey()
	clientAddress, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)

	msg = types.MsgOpenContract{
		Provider:     providerPubKey.String(),
		Service:      service.String(),
		Creator:      clientAddress.String(),
		Client:       clientPubKey.String(),
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         cosmos.NewInt64Coin("uarkeo", 15),
		Deposit:      cosmos.NewInt(1000),
	}
	require.NoError(t, k.MintAndSendToAccount(ctx, clientAddress, getCoin(common.Tokens(10))))
	require.NoError(t, s.OpenContractHandle(ctx, &msg))

	contract, err = k.GetActiveContractForUser(ctx, clientPubKey, providerPubKey, service)
	require.NoError(t, err)

	require.False(t, contract.IsEmpty())
	require.Equal(t, contract.Id, uint64(2))

	_, err = s.OpenContract(ctx, &msg)
	require.ErrorIs(t, err, types.ErrOpenContractAlreadyOpen)

	// confirm that the client can open a contract with a deleagate
	delegatePubKey := types.GetRandomPubKey()
	msg.Delegate = delegatePubKey.String()
	_, err = s.OpenContract(ctx, &msg)
	require.NoError(t, err)

	contract, err = k.GetActiveContractForUser(ctx, delegatePubKey, providerPubKey, service)
	require.NoError(t, err)

	require.False(t, contract.IsEmpty())
	require.Equal(t, contract.Id, uint64(3))
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
		SettlementDuration:  10,
	}
	err = s.ModProviderHandle(ctx, &modProviderMsg)

	require.NoError(t, err)

	interfaceRegistry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(interfaceRegistry)
	module.NewBasicManager().RegisterInterfaces(interfaceRegistry)
	types.RegisterInterfaces(interfaceRegistry)
	cdc := codec.NewProtoCodec(interfaceRegistry)

	kb := cKeys.NewInMemory(cdc)
	info, _, err := kb.NewMnemonic("whatever", cKeys.English, `m/44'/931'/0'/0/0`, "", hd.Secp256k1)
	require.NoError(t, err)
	pk, err := info.GetPubKey()
	require.NoError(t, err)
	clientPubKey, err := common.NewPubKeyFromCrypto(pk)
	require.NoError(t, err)
	clientAddress, err := clientPubKey.GetMyAddress()
	require.NoError(t, err)
	require.NoError(t, k.MintAndSendToAccount(ctx, clientAddress, getCoin(common.Tokens(10))))

	msg := types.MsgOpenContract{
		Provider:     providerPubKey.String(),
		Service:      service.String(),
		Creator:      clientAddress.String(),
		Client:       clientPubKey.String(),
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         cosmos.NewInt64Coin("uarkeo", 15),
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
	require.Equal(t, contract.Id, uint64(1))

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
		Nonce:      20,
	}
	message := claimMsg.GetBytesToSign("arkeo")
	claimMsg.Signature, _, err = kb.Sign("whatever", message, signing.SignMode_SIGN_MODE_DIRECT)
	require.NoError(t, err)
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
