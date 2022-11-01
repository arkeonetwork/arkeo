package types

import (
	"mercury/common"

	. "gopkg.in/check.v1"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgClaimContractIncomeSuite struct{}

var _ = Suite(&MsgClaimContractIncomeSuite{})

func (MsgClaimContractIncomeSuite) TestValidateBasic(c *C) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)

	// invalid address
	msg := MsgClaimContractIncome{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, sdkerrors.ErrInvalidAddress)

	msg = MsgClaimContractIncome{
		Creator: acct.String(),
		PubKey:  pubkey,
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, ErrInvalidChain)

	msg.Client = acct.String()
	msg.Chain = common.BTCChain
	msg.Height = 100
	err = msg.ValidateBasic()
	c.Assert(err, IsNil)

	// check bad client
	msg.Client = "bogus"
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, sdkerrors.ErrInvalidAddress)

	// check auth to cancel a specific contract
	msg = MsgClaimContractIncome{
		Creator: GetRandomBech32Addr().String(),
		PubKey:  pubkey,
		Client:  GetRandomBech32Addr().String(),
		Chain:   common.BTCChain,
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, ErrProviderBadSigner)
}
