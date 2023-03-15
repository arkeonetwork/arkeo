package types

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	// invalid address
	msg := MsgBondProvider{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	msg = MsgBondProvider{
		Creator:  acct.String(),
		Provider: pubkey,
	}
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrInvalidService)

	msg.Service = common.BTCService.String()
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, ErrInvalidBond)

	msg.Bond = cosmos.NewInt(500)
	err = msg.ValidateBasic()
	require.NoError(t, err)
}
