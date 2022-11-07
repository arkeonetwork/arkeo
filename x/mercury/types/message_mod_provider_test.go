package types

import (
	"mercury/common"

	. "gopkg.in/check.v1"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgModProviderSuite struct{}

var _ = Suite(&MsgModProviderSuite{})

func (MsgModProviderSuite) TestValidateBasic(c *C) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)

	// invalid address
	msg := MsgModProvider{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, sdkerrors.ErrInvalidAddress)

	// happy path
	msg = MsgModProvider{
		Creator:             acct.String(),
		PubKey:              pubkey,
		Chain:               common.BTCChain.String(),
		MinContractDuration: 12,
		MaxContractDuration: 30,
		MetadataURI:         "http://mad.hatter.net/test?foo=baz",
	}
	err = msg.ValidateBasic()
	c.Assert(err, IsNil)

	// URI is too long
	msg.MetadataURI = "http://mad.hatter.net/testsdkfjlsdkfjlsdfjsldfjkdsljflsdjfkdsjflsdjkfsdjlfsdjkfldsjflksjdfljsdlkfjsdlkfjdsklfjsdlkfjsdkljflksdjfklsdjflskdjflksdjflksdjfldsjflksdjfldskjflsdkfjsdlkjfksdljflskdjfsdlkjfdksljflsdkjfkldsjfsdlkfjlksdjfklsdjflkdsjfklsdjfsdkljflksdjflksdfjdklsjfl?foo=baz"
	err = msg.ValidateBasic()
	c.Check(err, ErrIs, ErrInvalidModProviderMetdataURI)
}
