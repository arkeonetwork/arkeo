package types

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/testutil/sample"
)

func TestMsgSetVersion_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgSetVersion
		err  error
	}{
		{
			name: "valid address",
			msg: MsgSetVersion{
				Creator: sample.AccAddress().String(),
				Version: 15,
			},
		},
		{
			name: "invalid version",
			msg: MsgSetVersion{
				Creator: sample.AccAddress().String(),
				Version: 0,
			},
			err: ErrInvalidVersion,
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
