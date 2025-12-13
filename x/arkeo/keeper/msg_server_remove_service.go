package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (k msgServer) RemoveService(goCtx context.Context, msg *types.MsgRemoveService) (*types.MsgRemoveServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.RemoveServiceValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.RemoveServiceHandle(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgRemoveServiceResponse{}, nil
}

func (k msgServer) RemoveServiceValidate(ctx cosmos.Context, msg *types.MsgRemoveService) error {
	if !types.IsAuthorityAllowed(k.GetAuthority(), msg.Creator) {
		return sdkerrors.ErrUnauthorized
	}
	// must exist
	if _, exists := k.GetService(ctx, msg.Name); !exists {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "service name not found: %s", msg.Name)
	}
	return msg.ValidateBasic()
}

func (k msgServer) RemoveServiceHandle(ctx cosmos.Context, msg *types.MsgRemoveService) error {
	name := strings.ToLower(msg.Name)
	if err := k.Keeper.RemoveService(ctx, name); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeRemoveService,
			sdk.NewAttribute("name", name),
		),
	)

	return nil
}
