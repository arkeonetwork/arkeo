package types

import (
	"testing"

	"github.com/arkeonetwork/arkeo/testutil/sample"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"
)

func TestMsgSetVersion_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetVersion
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgSetVersion{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgSetVersion{
				Creator: sample.AccAddress().String(),
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
