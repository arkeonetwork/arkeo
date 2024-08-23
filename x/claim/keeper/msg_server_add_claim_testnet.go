//go:build testnet

package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func (k msgServer) AddClaim(goCtx context.Context, msg *types.MsgAddClaim) (*types.MsgAddClaimResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	ctx.Logger().Info("receive add-claim request", "chain", msg.Chain.String(),
		"address", msg.Address,
		"creator", msg.Creator,
		"amount", msg.Amount)
	// This method is only provide on testnet for test purpose , so allow to override the record
	coin := sdk.NewCoin(types.DefaultClaimDenom, cosmos.NewInt(msg.Amount))
	claim := types.ClaimRecord{
		Chain:          msg.Chain,
		Address:        msg.Address,
		AmountClaim:    coin,
		AmountVote:     coin,
		AmountDelegate: coin,
		IsTransferable: false,
	}
	if msg.Chain == types.ARKEO {
		claim.IsTransferable = true
	}
	if err := k.Keeper.SetClaimRecord(ctx, claim); err != nil {
		return nil, fmt.Errorf("fail to save claim record to db,err: %w", err)
	}
	return &types.MsgAddClaimResponse{}, nil
}
