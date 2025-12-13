package keeper

import (
	"context"
	"fmt"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (k msgServer) UpdateService(goCtx context.Context, msg *types.MsgUpdateService) (*types.MsgUpdateServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.UpdateServiceValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.UpdateServiceHandle(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgUpdateServiceResponse{}, nil
}

func (k msgServer) UpdateServiceValidate(ctx cosmos.Context, msg *types.MsgUpdateService) error {
	if !types.IsAuthorityAllowed(k.GetAuthority(), msg.Creator) {
		return sdkerrors.ErrUnauthorized
	}

	// existing service must exist
	current, exists := k.GetService(ctx, msg.Name)
	if !exists {
		return errorsmod.Wrapf(sdkerrors.ErrNotFound, "service name not found: %s", msg.Name)
	}
	// id must match existing
	if current.Id != msg.Id {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "id mismatch for service %s: have %d, got %d", msg.Name, current.Id, msg.Id)
	}

	return msg.ValidateBasic()
}

func (k msgServer) UpdateServiceHandle(ctx cosmos.Context, msg *types.MsgUpdateService) error {
	svc := types.Service{
		Id:          msg.Id,
		Name:        strings.ToLower(msg.Name),
		Description: msg.Description,
		ServiceType: msg.ServiceType,
	}
	if err := k.SetService(ctx, svc); err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	sdkCtx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeUpdateService,
			sdk.NewAttribute("id", fmt.Sprintf("%d", svc.Id)),
			sdk.NewAttribute("name", svc.Name),
			sdk.NewAttribute("description", svc.Description),
			sdk.NewAttribute("type", svc.ServiceType),
		),
	)
	return nil
}
