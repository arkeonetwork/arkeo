package types

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/stretchr/testify/require"
)

func TestModProviderValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	coinRates, err := cosmos.ParseCoins("15uarkeo")
	require.NoError(t, err)
	rates := []*ContractRate{
		{
			MeterType: MeterType_PAY_PER_BLOCK,
			UserType:  UserType_SINGLE_USER,
			Rates:     coinRates,
		},
		{
			MeterType: MeterType_PAY_PER_CALL,
			UserType:  UserType_SINGLE_USER,
			Rates:     coinRates,
		},
	}

	// happy path
	msg := MsgModProvider{
		Creator:             acct,
		Provider:            pubkey,
		Service:             common.BTCService.String(),
		MinContractDuration: 12,
		MaxContractDuration: 30,
		MetadataUri:         "http://mad.hatter.net/test?foo=baz",
		Rates:               rates,
	}
	err = msg.ValidateBasic()
	require.NoError(t, err)

	// URI is too long
	msg.MetadataUri = "http://mad.hatter.net/testsdkfjlsdkfjlsdfjsldfjkdsljflsdjfkdsjflsdjkfsdjlfsdjkfldsjflksjdfljsdlkfjsdlkfjdsklfjsdlkfjsdkljflksdjfklsdjflskdjflksdjflksdjfldsjflksdjfldskjflsdkfjsdlkjfksdljflskdjfsdlkjfdksljflsdkjfkldsjfsdlkfjlksdjfklsdjflkdsjfklsdjfsdkljflksdjflksdfjdklsjfl?foo=baz"
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrInvalidModProviderMetdataURI)
}
