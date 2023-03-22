package types

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestModProviderValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require := require.New(t)
	require.NoError(err)

	// invalid address
	msg := MsgModProvider{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	require.ErrorIs(err, sdkerrors.ErrInvalidAddress)

	// happy path
	msg = MsgModProvider{
		Creator:             acct.String(),
		Provider:            pubkey,
		Service:             common.BTCService.String(),
		MinContractDuration: 12,
		MaxContractDuration: 30,
		MetadataUri:         "http://mad.hatter.net/test?foo=baz",
	}
	err = msg.ValidateBasic()
	require.NoError(err)

	// URI is too long
	msg.MetadataUri = "http://mad.hatter.net/testsdkfjlsdkfjlsdfjsldfjkdsljflsdjfkdsjflsdjkfsdjlfsdjkfldsjflksjdfljsdlkfjsdlkfjdsklfjsdlkfjsdkljflksdjfklsdjflskdjflksdjflksdjfldsjflksdjfldskjflsdkfjsdlkjfksdljflskdjfsdlkjfdksljflsdkjfkldsjfsdlkfjlksdjfklsdjflkdsjfklsdjfsdkljflksdjflksdfjdklsjfl?foo=baz"
	err = msg.ValidateBasic()
	require.ErrorIs(err, ErrInvalidModProviderMetdataURI)

	// check for duplicated rates.
	msg.Rates = []*ContractRate{
		{
			MeterType: MeterType_PAY_PER_BLOCK,
			UserType:  UserType_SINGLE_USER,
			Rate:      100,
		},
		{
			MeterType: MeterType_PAY_PER_BLOCK,
			UserType:  UserType_SINGLE_USER,
			Rate:      200,
		},
	}
	msg.MetadataUri = "http://mad.hatter.net/test?foo=baz"
	err = msg.ValidateBasic()
	require.ErrorIs(err, ErrInvalidModProviderDuplicateContractRates)

	msg.Rates = []*ContractRate{
		{
			MeterType: MeterType_PAY_PER_BLOCK,
			UserType:  UserType_SINGLE_USER,
			Rate:      100,
		},
		{
			MeterType: MeterType_PAY_PER_CALL,
			UserType:  UserType_SINGLE_USER,
			Rate:      200,
		},
	}
	err = msg.ValidateBasic()
	require.NoError(err)
}
