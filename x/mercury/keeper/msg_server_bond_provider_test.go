package keeper

import (
	"mercury/common"
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	. "gopkg.in/check.v1"
)

type BondProviderSuite struct{}

var _ = Suite(&BondProviderSuite{})

func (BondProviderSuite) TestValidate(c *C) {
}

func (BondProviderSuite) TestHandle(c *C) {
	ctx, k := SetupKeeper(c)

	s := newMsgServer(k)

	// setup
	pubkey := types.GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	c.Assert(k.MintAndSendToAccount(ctx, acct, getCoin(tokens(10))), IsNil)

	// Add to bond
	msg := types.MsgBondProvider{
		Creator: acct.String(),
		Pubkey:  pubkey,
		Chain:   common.BTCChain,
		Bond:    cosmos.NewInt(tokens(8)),
	}
	c.Assert(s.BondProviderHandle(ctx, &msg), IsNil)
	// check balance as drawn down by two
	bal := k.GetBalance(ctx, acct)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, tokens(2))
	// check that provider now exists
	c.Check(k.ProviderExists(ctx, msg.Pubkey, msg.Chain), Equals, true)
	provider, err := k.GetProvider(ctx, msg.Pubkey, msg.Chain)
	c.Assert(err, IsNil)
	c.Check(provider.Bond.Int64(), Equals, tokens(8))

	// remove too much bond
	msg.Bond = cosmos.NewInt(tokens(-20))
	err = s.BondProviderHandle(ctx, &msg)
	c.Check(err, ErrIs, types.ErrInsufficientFunds)

	// check balance hasn't changed
	bal = k.GetBalance(ctx, acct)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, tokens(2))
	// check provider has same bond
	provider, err = k.GetProvider(ctx, msg.Pubkey, msg.Chain)
	c.Assert(err, IsNil)
	c.Check(provider.Bond.Int64(), Equals, tokens(8))

	// remove all bond
	msg.Bond = cosmos.NewInt(tokens(-8))
	err = s.BondProviderHandle(ctx, &msg)
	c.Assert(err, IsNil)

	bal = k.GetBalance(ctx, acct) // check balance
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, tokens(10))
	c.Check(k.ProviderExists(ctx, msg.Pubkey, msg.Chain), Equals, false) // should be removed
}
