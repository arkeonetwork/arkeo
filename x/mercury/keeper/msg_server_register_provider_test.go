package keeper

import (
	"mercury/common"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	. "gopkg.in/check.v1"
)

type RegisterProviderSuite struct{}

var _ = Suite(&RegisterProviderSuite{})

func (RegisterProviderSuite) TestValidate(c *C) {
	ctx, k := SetupKeeper(c)

	s := newMsgServer(k)

	// setup
	pubkey := types.GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	c.Assert(k.MintAndSendToAccount(ctx, acct, getCoin(10_00000000)), IsNil)

	msg := types.MsgRegisterProvider{
		Creator: acct.String(),
		Pubkey:  pubkey,
		Chain:   common.BTCChain,
	}

	// happy path
	c.Assert(s.RegisterProviderValidate(ctx, &msg), IsNil)

	// mismatch of address and pubkey
	badAcct := types.GetRandomBech32Addr()
	msg.Creator = badAcct.String()
	err = s.RegisterProviderValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrProviderBadSigner)

	// provider already exists
	msg.Creator = acct.String()
	provider := types.NewProvider(msg.Pubkey, msg.Chain)
	c.Assert(k.SetProvider(ctx, provider), IsNil)
	err = s.RegisterProviderValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrProviderAlreadyExists)
	k.RemoveProvider(ctx, msg.Pubkey, msg.Chain)

	// insufficient funds
	c.Assert(k.SendFromAccountToModule(ctx, acct, types.ModuleName, getCoins(9_00000000)), IsNil)
	err = s.RegisterProviderValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrInsufficientFunds)
}

func (RegisterProviderSuite) TestHandle(c *C) {
	ctx, k := SetupKeeper(c)

	s := newMsgServer(k)

	// setup
	pubkey := types.GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	c.Assert(k.MintAndSendToAccount(ctx, acct, getCoin(10_00000000)), IsNil)

	msg := types.MsgRegisterProvider{
		Creator: acct.String(),
		Pubkey:  pubkey,
		Chain:   common.BTCChain,
	}

	c.Assert(s.RegisterProviderHandle(ctx, &msg), IsNil)

	// check balance as drawn down by two
	bal := k.GetBalance(ctx, acct)
	c.Check(bal.AmountOf(configs.Denom).Int64(), Equals, int64(8_00000000))

	// check that provider now exists
	c.Check(k.ProviderExists(ctx, msg.Pubkey, msg.Chain), Equals, true)
}
