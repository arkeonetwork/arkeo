package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/testutil/sample"
)

func TestMsgAddClaim_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgAddClaim
		err  error
	}{
		{
			name: "invalid chain",
			msg: MsgAddClaim{
				Creator: sample.AccAddress(),
				Address: sample.AccAddress().String(),
				Chain:   Chain(100),
				Amount:  100,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid amount",
			msg: MsgAddClaim{
				Creator: sample.AccAddress(),
				Address: sample.AccAddress().String(),
				Chain:   ARKEO,
				Amount:  -100,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: MsgAddClaim{
				Creator: sample.AccAddress(),
				Address: sample.AccAddress().String(),
				Chain:   ARKEO,
				Amount:  100,
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
