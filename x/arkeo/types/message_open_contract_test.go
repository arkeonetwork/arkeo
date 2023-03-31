package types

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/stretchr/testify/require"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TestBondProviderValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	// invalid address
	msg := MsgOpenContract{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	rate, err := cosmos.ParseCoin("0uarkeo")
	require.NoError(t, err)

	msg = MsgOpenContract{
		Creator:  acct.String(),
		Provider: pubkey,
		Client:   pubkey,
		Service:  common.BTCService.String(),
		Rate:     rate,
	}
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrOpenContractDuration)

	msg.Duration = 100
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrOpenContractRate)

	msg.Rate, _ = cosmos.ParseCoin("100uarkeo")
	err = msg.ValidateBasic()
	require.NoError(t, err)
}
