package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
)

func TestValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	msg := MsgBondProvider{
		Creator:  acct.String(),
		Provider: pubkey.String(),
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
