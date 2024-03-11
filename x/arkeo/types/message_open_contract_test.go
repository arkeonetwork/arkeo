package types

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/stretchr/testify/require"
)

func TestBondProviderValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	rate, err := cosmos.ParseCoin("0uarkeo")
	require.NoError(t, err)

	msg := MsgOpenContract{
		Creator:          acct.String(),
		Provider:         pubkey.String(),
		Client:           pubkey.String(),
		Service:          common.BTCService.String(),
		Rate:             rate,
		QueriesPerMinute: 10,
	}
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrOpenContractDuration)

	msg.Duration = 100
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrOpenContractRate)

	msg.Rate, _ = cosmos.ParseCoin("100uarkeo")
	err = msg.ValidateBasic()
	require.NoError(t, err)

	msg.Authorization = ContractAuthorization_OPEN
	msg.ContractType = ContractType_PAY_AS_YOU_GO
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrInvalidAuthorization)
}
