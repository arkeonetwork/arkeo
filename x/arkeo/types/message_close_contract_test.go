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
		Creator:    acct,
		ContractId: 50,
	}
	err = msg.ValidateBasic()
	require.NoError(t, err)
}
