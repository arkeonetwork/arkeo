//go:build !testnet

package keeper

import (
	"context"
	"fmt"

	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func (k msgServer) AddClaim(goCtx context.Context, msg *types.MsgAddClaim) (*types.MsgAddClaimResponse, error) {
	return nil, fmt.Errorf("MsgAddClaim is only support on testnet")
}
