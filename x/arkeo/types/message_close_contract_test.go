package types

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCloseContractValidateBasic(t *testing.T) {
	// setup
	pubkey := GetRandomPubKey()
	acct, err := pubkey.GetMyAddress()
	require.NoError(t, err)

	msg := MsgCloseContract{
		Creator:    acct.String(),
		ContractId: 50,
		Client:     pubkey,
	}
	err = msg.ValidateBasic()
	require.NoError(t, err)
}
