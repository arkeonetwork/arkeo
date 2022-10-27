package types

import (
	"mercury/common"

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
		Creator: acct.String(),
		PubKey:  pubkey,
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, ErrInvalidChain)

	msg.Chain = common.BTCChain
	err = msg.ValidateBasic()
	c.Assert(err, IsNil)

	// check auth to cancel a specific contract
	msg = MsgCloseContract{
		Creator: GetRandomBech32Addr().String(),
		PubKey:  pubkey,
		Client:  GetRandomBech32Addr().String(),
		Chain:   common.BTCChain,
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, ErrProviderBadSigner)

	msg.Client = "bogus"
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, sdkerrors.ErrInvalidAddress)
}
