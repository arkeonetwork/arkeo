package types

import (
	"mercury/common/cosmos"

	. "gopkg.in/check.v1"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgBondProviderSuite struct{}

var _ = Suite(&MsgBondProviderSuite{})

func (MsgBondProviderSuite) TestValidateBasic(c *C) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)

	// invalid address
	msg := MsgBondProvider{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, sdkerrors.ErrInvalidAddress)

	msg = MsgBondProvider{
		Creator: acct.String(),
		Pubkey:  pubkey,
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, ErrInvalidBond)

	msg.Bond = cosmos.NewInt(500)
	err = msg.ValidateBasic()
	c.Assert(err, IsNil)
}
