package keeper

import (
	"mercury/common"
	"mercury/common/cosmos"
	"mercury/x/mercury/types"

	. "gopkg.in/check.v1"
)

type OpenContractSuite struct{}

var _ = Suite(&OpenContractSuite{})

func (OpenContractSuite) TestValidate(c *C) {
	var err error
	ctx, k := SetupKeeper(c)

	s := newMsgServer(k)

	// setup
	pubkey := types.GetRandomPubKey()
	acc := types.GetRandomBech32Addr()
	chain := common.BTCChain

	provider := types.NewProvider(pubkey, chain)
	provider.Bond = cosmos.NewInt(500_00000000)
	provider.MaxContractDuration = 1000
	provider.MinContractDuration = 10
	provider.SubscriptionRate = 15
	provider.PayAsYouGoRate = 2
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	// happy path
	msg := types.MsgOpenContract{
		PubKey:   pubkey,
		Chain:    chain,
		Creator:  acc.String(),
		CType:    types.ContractType_Subscription,
		Duration: 100,
		Rate:     15,
	}
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
	msg.CType = types.ContractType_PayAsYouGo
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrOpenContractMismatchRate)
	msg.Rate = 15
	msg.CType = types.ContractType_Subscription

	provider.Bond = cosmos.NewInt(1)
	c.Assert(k.SetProvider(ctx, provider), IsNil)
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrInvalidBond)
	provider.Bond = cosmos.NewInt(500_00000000)
	c.Assert(k.SetProvider(ctx, provider), IsNil)

	ctx = ctx.WithBlockHeight(14)
	contract := types.NewContract(pubkey, chain, acc)
	contract.Type = types.ContractType_Subscription
	contract.Height = ctx.BlockHeight()
	contract.Duration = 100
	contract.Rate = 2
	c.Assert(k.SetContract(ctx, contract), IsNil)
	err = s.OpenContractValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrOpenContractAlreadyOpen)
}

func (OpenContractSuite) TestHandle(c *C) {
	ctx, k := SetupKeeper(c)

	s := newMsgServer(k)

	// setup
	pubkey := types.GetRandomPubKey()
	acc, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	chain := common.BTCChain

	msg := types.MsgOpenContract{
		PubKey:   pubkey,
		Chain:    chain,
		Creator:  acc.String(),
		CType:    types.ContractType_PayAsYouGo,
		Duration: 100,
		Rate:     15,
	}
	c.Assert(s.OpenContractHandle(ctx, &msg), IsNil)

	contract, err := k.GetContract(ctx, pubkey, chain, acc)
	c.Assert(err, IsNil)

	c.Check(contract.Type, Equals, types.ContractType_PayAsYouGo)
	c.Check(contract.Height, Equals, ctx.BlockHeight())
	c.Check(contract.Duration, Equals, int64(100))
	c.Check(contract.Rate, Equals, int64(15))
}
