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
			name: "invalid creator address",
			msg: MsgAddClaim{
				Creator: "invalid_address",
				Address: sample.AccAddress(),
				Chain:   "swapi.dev",
				Amount:  100,
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "invalid address",
			msg: MsgAddClaim{
				Creator: sample.AccAddress(),
				Address: "invalid address",
				Chain:   "swapi.dev",
				Amount:  100,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid chain",
			msg: MsgAddClaim{
				Creator: sample.AccAddress(),
				Address: sample.AccAddress(),
				Chain:   "invalid chain",
				Amount:  100,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "invalid amount",
			msg: MsgAddClaim{
				Creator: sample.AccAddress(),
				Address: sample.AccAddress(),
				Chain:   "swapi.dev",
				Amount:  -100,
			},
			err: sdkerrors.ErrInvalidRequest,
		},
		{
			name: "valid address",
			msg: MsgAddClaim{
				Creator: sample.AccAddress(),
				Address: sample.AccAddress(),
				Chain:   "swapi.dev",
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
