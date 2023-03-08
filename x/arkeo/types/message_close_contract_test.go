package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestCloseContractValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	// invalid address
	msg := MsgCloseContract{
		Creator: "invalid address",
	}
	err = msg.ValidateBasic()
	require.ErrorIs(t, err, sdkerrors.ErrInvalidAddress)

	msg = MsgCloseContract{
		Creator:    acct.String(),
		ContractId: 50,
	}
	err = msg.ValidateBasic()
	require.NoError(t, err)
}
