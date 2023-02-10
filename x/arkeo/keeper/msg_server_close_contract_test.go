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
	pubkey := types.GetRandomPubKey()
	acc := types.GetRandomPubKey()
	chain := common.BTCChain

	contract := types.NewContract(pubkey, chain, acc)
	contract.Duration = 100
	contract.Height = 10
	c.Assert(k.SetContract(ctx, contract), IsNil)

	// happy path
	msg := types.MsgCloseContract{
		PubKey:  pubkey,
		Chain:   chain.String(),
		Creator: acc.String(),
		Client:  acc,
	}
	c.Assert(s.CloseContractValidate(ctx, &msg), IsNil)

	contract.Duration = 3
	c.Assert(k.SetContract(ctx, contract), IsNil)
	err := s.CloseContractValidate(ctx, &msg)
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
	contract.Type = types.ContractType_Subscription
	contract.Duration = 100
	contract.Height = 10
	contract.Rate = 5
	contract.Deposit = cosmos.NewInt(500)
	c.Assert(k.SetContract(ctx, contract), IsNil)

	// happy path
	msg := types.MsgCloseContract{
		PubKey:  pubkey,
		Chain:   chain.String(),
		Creator: acc.String(),
		Client:  acc,
	}
	c.Assert(s.CloseContractHandle(ctx, &msg), IsNil)

	contract, err = k.GetContract(ctx, pubkey, chain, acc)
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
