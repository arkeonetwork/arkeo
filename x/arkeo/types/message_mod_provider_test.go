package types

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestModProviderValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	// invalid address
	msg := MsgModProvider{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	rates, err := cosmos.ParseCoins("15uarkeo")
	require.NoError(t, err)

	// happy path
	msg = MsgModProvider{
		Creator:             acct.String(),
		Provider:            pubkey,
		Service:             common.BTCService.String(),
		MinContractDuration: 12,
		MaxContractDuration: 30,
		MetadataUri:         "http://mad.hatter.net/test?foo=baz",
		SubscriptionRate:    rates,
		PayAsYouGoRate:      rates,
	}
	err = msg.ValidateBasic()
	require.NoError(t, err)

	// URI is too long
	msg.MetadataUri = "http://mad.hatter.net/testsdkfjlsdkfjlsdfjsldfjkdsljflsdjfkdsjflsdjkfsdjlfsdjkfldsjflksjdfljsdlkfjsdlkfjdsklfjsdlkfjsdkljflksdjfklsdjflskdjflksdjflksdjfldsjflksdjfldskjflsdkfjsdlkjfksdljflskdjfsdlkjfdksljflsdkjfkldsjfsdlkfjlksdjfklsdjflkdsjfklsdjfsdkljflksdjflksdfjdklsjfl?foo=baz"
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrInvalidModProviderMetdataURI)
}
