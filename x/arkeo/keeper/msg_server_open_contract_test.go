package keeper

import (
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	. "gopkg.in/check.v1"
)

type OpenContractSuite struct{}

var _ = Suite(&OpenContractSuite{})

func (OpenContractSuite) TestValidate(c *C) {
	var err error
	ctx, k, sk := SetupKeeperWithStaking(c)
	ctx = ctx.WithBlockHeight(1)
	s := newMsgServer(k, sk)

	// setup
	providerPubkey := types.GetRandomPubKey()
	clientPubKey := types.GetRandomPubKey()
	acc, err := clientPubKey.GetMyAddress()
	if err != nil {
		c.Error(err)
	}

	chain := common.BTCChain

	provider := types.NewProvider(providerPubkey, chain)
	provider.Bond = cosmos.NewInt(500_00000000)
	provider.Status = types.ProviderStatus_ONLINE
	provider.MaxContractDuration = 1000
	provider.MinContractDuration = 10
	provider.SubscriptionRate = 15
	provider.PayAsYouGoRate = 2
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	// happy path
	msg := types.MsgOpenContract{
		PubKey:       providerPubkey,
		Chain:        chain.String(),
		Client:       clientPubKey,
		Creator:      acc.String(),
		ContractType: types.ContractType_SUBSCRIPTION,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(100 * 15),
	}
	c.Assert(k.MintAndSendToAccount(ctx, acc, getCoin(common.Tokens(100*25))), IsNil)
	c.Assert(s.OpenContractValidate(ctx, &msg), IsNil)

	// check duration
	msg.Duration = 10000000000000
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrOpenContractDuration)
	msg.Duration = 5
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrOpenContractDuration)
	msg.Duration = 100

	// check rates
	msg.Rate = 10
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrOpenContractMismatchRate)
	msg.ContractType = types.ContractType_PAY_AS_YOU_GO
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrOpenContractMismatchRate)
	msg.Rate = 15
	msg.ContractType = types.ContractType_SUBSCRIPTION

	provider.Bond = cosmos.NewInt(1)
	c.Assert(k.SetProvider(ctx, provider), IsNil)
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrInvalidBond)
	provider.Bond = cosmos.NewInt(500_00000000)
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	ctx = ctx.WithBlockHeight(15)
	contract := types.NewContract(providerPubkey, chain, clientPubKey)
	contract.Type = types.ContractType_SUBSCRIPTION
	contract.Height = ctx.BlockHeight()
	contract.Duration = 100
	contract.Rate = 2

	c.Assert(s.OpenContractHandle(ctx, &msg), IsNil)
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrOpenContractAlreadyOpen)
}

func (OpenContractSuite) TestHandle(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	// setup
	pubkey := types.GetRandomPubKey()
	acc, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	chain := common.BTCChain
	c.Assert(k.MintAndSendToAccount(ctx, acc, getCoin(common.Tokens(10))), IsNil)

	msg := types.MsgOpenContract{
		PubKey:       pubkey,
		Chain:        chain.String(),
		Creator:      acc.String(),
		Client:       pubkey,
		ContractType: types.ContractType_PAY_AS_YOU_GO,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1000),
	}
	c.Assert(s.OpenContractHandle(ctx, &msg), IsNil)

	contract, err := k.GetActiveContractForUser(ctx, pubkey, pubkey, chain)
	c.Assert(err, IsNil)

	c.Check(contract.Type, Equals, types.ContractType_PAY_AS_YOU_GO)
	c.Check(contract.IsEmpty(), Equals, false)

	c.Check(contract.Height, Equals, ctx.BlockHeight())
	c.Check(contract.Duration, Equals, int64(100))
	c.Check(contract.Rate, Equals, int64(15))
	c.Check(contract.Nonce, Equals, int64(0))
	c.Check(contract.Deposit.Int64(), Equals, int64(1000))
	c.Check(contract.Paid.Int64(), Equals, int64(0))

	bal := k.GetBalance(ctx, acc) // check balance
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(899999000))

	// check that contract expiration has been set
	set, err := k.GetContractExpirationSet(ctx, contract.Expiration())
	c.Assert(err, IsNil)
	c.Check(set.Height, Equals, contract.Expiration())
	c.Check(set.ContractSet.ContractIds, HasLen, 1)

	// check that contract has been added to the user
	userSet, err := k.GetUserContractSet(ctx, contract.GetSpender())
	c.Assert(err, IsNil)
	c.Check(userSet.User, Equals, contract.GetSpender())
	c.Check(userSet.ContractSet.ContractIds, HasLen, 1)
}

func (OpenContractSuite) TestOpenContract(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)
	ctx = ctx.WithBlockHeight(10)
	s := newMsgServer(k, sk)

	providerPubKey := types.GetRandomPubKey()
	providerAddress, err := providerPubKey.GetMyAddress()
	c.Assert(err, IsNil)
	chain := common.BTCChain
	provider := types.NewProvider(providerPubKey, chain)
	provider.Bond = cosmos.NewInt(10000000000)
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	modProviderMsg := types.MsgModProvider{
		PubKey:              provider.PubKey,
		Chain:               provider.Chain.String(),
		MinContractDuration: 10,
		MaxContractDuration: 500,
		Status:              types.ProviderStatus_Online,
		PayAsYouGoRate:      15,
		SubscriptionRate:    15,
	}
	err = s.ModProviderHandle(ctx, &modProviderMsg)

	c.Assert(err, IsNil)
	c.Assert(k.MintAndSendToAccount(ctx, providerAddress, getCoin(common.Tokens(10))), IsNil)

	msg := types.MsgOpenContract{
		PubKey:       providerPubKey,
		Chain:        chain.String(),
		Creator:      providerAddress.String(),
		Client:       providerPubKey,
		ContractType: types.ContractType_PayAsYouGo,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1500),
	}
	_, err = s.OpenContract(ctx, &msg)
	c.Assert(err, IsNil)

	contract, err := k.GetActiveContractForUser(ctx, providerPubKey, providerPubKey, chain)
	c.Assert(err, IsNil)

	c.Check(contract.IsEmpty(), Equals, false)
	c.Check(contract.Id, Equals, uint64(0))

	clientPubKey := types.GetRandomPubKey()
	clientAddress, err := clientPubKey.GetMyAddress()
	c.Assert(err, IsNil)

	msg = types.MsgOpenContract{
		PubKey:       providerPubKey,
		Chain:        chain.String(),
		Creator:      clientAddress.String(),
		Client:       clientPubKey,
		ContractType: types.ContractType_PayAsYouGo,
		Duration:     100,
		Rate:         15,
		Deposit:      cosmos.NewInt(1000),
	}
	c.Assert(k.MintAndSendToAccount(ctx, clientAddress, getCoin(common.Tokens(10))), IsNil)
	c.Assert(s.OpenContractHandle(ctx, &msg), IsNil)

	contract, err = k.GetActiveContractForUser(ctx, clientPubKey, providerPubKey, chain)
	c.Assert(err, IsNil)

	c.Check(contract.IsEmpty(), Equals, false)
	c.Check(contract.Id, Equals, uint64(1))

	_, err = s.OpenContract(ctx, &msg)
	c.Check(err, ErrIs, types.ErrOpenContractAlreadyOpen)
}
