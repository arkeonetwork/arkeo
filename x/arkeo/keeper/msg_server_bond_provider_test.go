package keeper

import (
	"github.com/ArkeoNetwork/arkeo-protocol/common"
	"github.com/ArkeoNetwork/arkeo-protocol/common/cosmos"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/configs"
	"github.com/ArkeoNetwork/arkeo-protocol/x/arkeo/types"

	. "gopkg.in/check.v1"
)

type BondProviderSuite struct{}

var _ = Suite(&BondProviderSuite{})

func (BondProviderSuite) TestValidate(c *C) {
}

func (BondProviderSuite) TestHandle(c *C) {
	ctx, k, sk := SetupKeeperWithStaking(c)

	s := newMsgServer(k, sk)

	// setup
	pubkey := types.GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	c.Assert(k.MintAndSendToAccount(ctx, acct, getCoin(common.Tokens(10))), IsNil)

	// Add to bond
	msg := types.MsgBondProvider{
		Creator: acct.String(),
		PubKey:  pubkey,
		Chain:   common.BTCChain.String(),
		Bond:    cosmos.NewInt(common.Tokens(8)),
	}
	c.Assert(s.BondProviderHandle(ctx, &msg), IsNil)
	// check balance as drawn down by two
	bal := k.GetBalance(ctx, acct)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, common.Tokens(2))
	// check that provider now exists
	c.Check(k.ProviderExists(ctx, msg.PubKey, common.BTCChain), Equals, true)
	provider, err := k.GetProvider(ctx, msg.PubKey, common.BTCChain)
	c.Assert(err, IsNil)
	c.Check(provider.Bond.Int64(), Equals, common.Tokens(8))

	// remove too much bond
	msg.Bond = cosmos.NewInt(common.Tokens(-20))
	err = s.BondProviderHandle(ctx, &msg)
	c.Check(err, ErrIs, types.ErrInsufficientFunds)

	// check balance hasn't changed
	bal = k.GetBalance(ctx, acct)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, common.Tokens(2))
	// check provider has same bond
	provider, err = k.GetProvider(ctx, msg.PubKey, common.BTCChain)
	c.Assert(err, IsNil)
	c.Check(provider.Bond.Int64(), Equals, common.Tokens(8))

	// remove all bond
	msg.Bond = cosmos.NewInt(common.Tokens(-8))
	err = s.BondProviderHandle(ctx, &msg)
	c.Assert(err, IsNil)

	bal = k.GetBalance(ctx, acct) // check balance
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, common.Tokens(10))
	c.Check(k.ProviderExists(ctx, msg.PubKey, common.BTCChain), Equals, false) // should be removed
}
