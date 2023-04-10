package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/testutil/sample"
)

func TestMsgTransferClaim_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTransferClaim
		err  error
	}{
		{
			name: "valid create address and valid to address",
			msg: MsgTransferClaim{
				Creator:   sample.AccAddress(),
				ToAddress: sample.AccAddress(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
