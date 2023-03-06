package types

import (
	. "gopkg.in/check.v1"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgCloseContractSuite struct{}

var _ = Suite(&MsgCloseContractSuite{})

func (MsgCloseContractSuite) TestValidateBasic(c *C) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)

	// invalid address
	msg := MsgCloseContract{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, sdkerrors.ErrInvalidAddress)

	msg = MsgCloseContract{
		Creator:    acct.String(),
		ContractId: 50,
	}
	err = msg.ValidateBasic()
	c.Assert(err, IsNil)
}
