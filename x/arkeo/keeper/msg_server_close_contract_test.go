package keeper

import (
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	. "gopkg.in/check.v1"
)

type CloseContractSuite struct{}

var _ = Suite(&CloseContractSuite{})

func (CloseContractSuite) TestValidate(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)
	ctx = ctx.WithBlockHeight(14)

	s := newMsgServer(k, sk)

	// setup
	providerPubkey := types.GetRandomPubKey()

	clientPubKey := types.GetRandomPubKey()
	clientAcct, err := clientPubKey.GetMyAddress()
	c.Assert(err, IsNil)
	chain := common.BTCChain

	contract := types.NewContract(providerPubkey, chain, clientPubKey)
	contract.Duration = 100
	contract.Height = 10
	contract.Id = 1
	c.Assert(k.SetContract(ctx, contract), IsNil)

	// happy path
	msg := types.MsgCloseContract{
		Creator:    clientAcct.String(),
		ContractId: contract.Id,
	}
	c.Assert(s.CloseContractValidate(ctx, &msg), IsNil)

	contract.Duration = 3
	c.Assert(k.SetContract(ctx, contract), IsNil)
	err = s.CloseContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrCloseContractAlreadyClosed)
}

func (CloseContractSuite) TestHandle(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)
	ctx = ctx.WithBlockHeight(14)

	s := newMsgServer(k, sk)

	// setup
	c.Assert(k.MintToModule(ctx, types.ModuleName, getCoin(500)), IsNil)
	c.Assert(k.SendFromModuleToModule(ctx, types.ModuleName, types.ContractName, getCoins(500)), IsNil)
	pubkey := types.GetRandomPubKey()
	provider, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	acc := types.GetRandomPubKey()
	chain := common.BTCChain
	c.Check(k.GetBalance(ctx, provider).IsZero(), Equals, true)

	contract := types.NewContract(pubkey, chain, acc)
	contract.Type = types.ContractType_SUBSCRIPTION
	contract.Duration = 100
	contract.Height = 10
	contract.Rate = 5
	contract.Deposit = cosmos.NewInt(500)
	contract.Id = 1
	c.Assert(k.SetContract(ctx, contract), IsNil)

	// happy path
	msg := types.MsgCloseContract{
		Creator:    acc.String(),
		ContractId: contract.Id,
	}
	c.Assert(s.CloseContractHandle(ctx, &msg), IsNil)

	contract, err = k.GetContract(ctx, contract.Id)
	c.Assert(err, IsNil)
	c.Check(contract.Paid.Int64(), Equals, int64(20))
	c.Check(contract.ClosedHeight, Equals, ctx.BlockHeight())

	bal := k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom)
	c.Check(bal.Int64(), Equals, int64(0))
	c.Check(k.HasCoins(ctx, provider, getCoins(18)), Equals, true)
	c.Check(k.HasCoins(ctx, contract.ClientAddress(), getCoins(480)), Equals, true)
	bal = k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom)
	c.Check(bal.Int64(), Equals, int64(2))
}

func (CloseContractSuite) TestCloseSubscriptionContract(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	// set up provider
	providerPubKey := types.GetRandomPubKey()
	providerAddress, err := providerPubKey.GetMyAddress()
	c.Assert(err, IsNil)
	chain := common.BTCChain
	provider := types.NewProvider(providerPubKey, chain)
	provider.Bond = cosmos.NewInt(10000000000)
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	modProviderMsg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Chain:               provider.Chain.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      15,
		SubscriptionRate:    15,
	}
	err = s.ModProviderHandle(ctx, &modProviderMsg)

	c.Assert(err, IsNil)
	c.Assert(k.MintAndSendToAccount(ctx, providerAddress, getCoin(common.Tokens(10))), IsNil)

	// set up contract with no delegate address
	userPubKey := types.GetRandomPubKey()
	userAddress, err := userPubKey.GetMyAddress()
	c.Assert(err, IsNil)

	openContractMessage := types.MsgOpenContract{
		Provider:     providerPubKey,
		Chain:        chain.String(),
		Creator:      providerAddress.String(),
		Client:       userPubKey,
		ContractType: types.ContractType_SUBSCRIPTION,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &openContractMessage)
	c.Assert(err, IsNil)

	contract, err := s.GetActiveContractForUser(ctx, userPubKey, providerPubKey, chain)
	c.Assert(err, IsNil)

	// confirm that another user cannot close the contract
	use2PubKey := types.GetRandomPubKey()
	user2Address, err := use2PubKey.GetMyAddress()
	c.Assert(err, IsNil)

	closeContractMsg := types.MsgCloseContract{
		Creator:    user2Address.String(),
		ContractId: contract.Id,
	}
	_, err = s.CloseContract(ctx, &closeContractMsg)
	c.Check(err, ErrIs, types.ErrCloseContractUnauthorized)

	// confirm that the contract can be closed by the client
	closeContractMsg.Creator = userAddress.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	c.Check(err, IsNil)

	// reopen contract this time with a delagate address.
	openContractMessage.Delegate = use2PubKey
	_, err = s.OpenContract(ctx, &openContractMessage)
	c.Assert(err, IsNil)

	contract, err = s.GetActiveContractForUser(ctx, userPubKey, providerPubKey, chain)
	c.Assert(err, IsNil)
	closeContractMsg.ContractId = contract.Id

	// confirm that the contract cannot be closed by the delegate
	closeContractMsg.Creator = user2Address.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	c.Check(err, ErrIs, types.ErrCloseContractUnauthorized)

	// but can be closed by the client
	closeContractMsg.Creator = userAddress.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	c.Check(err, IsNil)
}

func (CloseContractSuite) TestClosePayAsYouGoContract(c *C) {
	// NOTE: pay as you go contracts cannot be closed on demand.
	ctx, k, sk := SetupKeeperWithStaking(c)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	// set up provider
	providerPubKey := types.GetRandomPubKey()
	providerAddress, err := providerPubKey.GetMyAddress()
	c.Assert(err, IsNil)
	chain := common.BTCChain
	provider := types.NewProvider(providerPubKey, chain)
	provider.Bond = cosmos.NewInt(10000000000)
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	modProviderMsg := types.MsgModProvider{
		Provider:            provider.PubKey,
		Chain:               provider.Chain.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_ONLINE,
		PayAsYouGoRate:      15,
		SubscriptionRate:    15,
	}
	err = s.ModProviderHandle(ctx, &modProviderMsg)

	c.Assert(err, IsNil)
	c.Assert(k.MintAndSendToAccount(ctx, providerAddress, getCoin(common.Tokens(10))), IsNil)

	// set up contract with no delegate address
	userPubKey := types.GetRandomPubKey()
	userAddress, err := userPubKey.GetMyAddress()
	c.Assert(err, IsNil)

	openContractMessage := types.MsgOpenContract{
		Provider:     providerPubKey,
		Chain:        chain.String(),
		Creator:      providerAddress.String(),
		Client:       userPubKey,
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &openContractMessage)
	c.Assert(err, IsNil)

	contract, err := s.GetActiveContractForUser(ctx, userPubKey, providerPubKey, chain)
	c.Assert(err, IsNil)

	// confirm that another user cannot close the contract
	use2PubKey := types.GetRandomPubKey()
	user2Address, err := use2PubKey.GetMyAddress()
	c.Assert(err, IsNil)

	closeContractMsg := types.MsgCloseContract{
		Creator:    user2Address.String(),
		ContractId: contract.Id,
	}
	_, err = s.CloseContract(ctx, &closeContractMsg)
	c.Check(err, ErrIs, types.ErrCloseContractUnauthorized)

	// confirm that the contract can not be closed by the client
	closeContractMsg.Creator = userAddress.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	c.Check(err, ErrIs, types.ErrCloseContractUnauthorized)

	// reopen contract this time with a delagate address.
	openContractMessage.Delegate = use2PubKey
	_, err = s.OpenContract(ctx, &openContractMessage)
	c.Assert(err, IsNil)

	contract, err = s.GetActiveContractForUser(ctx, userPubKey, providerPubKey, chain)
	c.Assert(err, IsNil)

	closeContractMsg.ContractId = contract.Id

	// confirm that the contract cannot be closed by the delegate
	closeContractMsg.Creator = user2Address.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	c.Check(err, ErrIs, types.ErrCloseContractUnauthorized)

	// nor the client
	closeContractMsg.Creator = userAddress.String()
	_, err = s.CloseContract(ctx, &closeContractMsg)
	c.Check(err, ErrIs, types.ErrCloseContractUnauthorized)
}
