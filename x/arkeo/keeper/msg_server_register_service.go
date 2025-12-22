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

func (k msgServer) RegisterService(goCtx context.Context, msg *types.MsgRegisterService) (*types.MsgRegisterServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if err := k.RegisterServiceValidate(ctx, msg); err != nil {
		return nil, err
	}

	if err := k.RegisterServiceHandle(ctx, msg); err != nil {
		return nil, err
	}

	return &types.MsgRegisterServiceResponse{}, nil
}

func (k msgServer) RegisterServiceValidate(ctx cosmos.Context, msg *types.MsgRegisterService) error {
	if !types.IsAuthorityAllowed(k.GetAuthority(), msg.Creator) {
		return sdkerrors.ErrUnauthorized
	}

	// Ensure unique name
	if _, exists := k.GetService(ctx, msg.Name); exists {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "service name already exists: %s", msg.Name)
	}
	// Ensure unique id
	if _, exists := k.GetServiceByID(ctx, msg.Id); exists {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "service id already exists: %d", msg.Id)
	}

	return msg.ValidateBasic()
}

func (k msgServer) RegisterServiceHandle(ctx cosmos.Context, msg *types.MsgRegisterService) error {
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
			types.EventTypeRegisterService,
			sdk.NewAttribute("id", fmt.Sprintf("%d", svc.Id)),
			sdk.NewAttribute("name", svc.Name),
			sdk.NewAttribute("description", svc.Description),
			sdk.NewAttribute("type", svc.ServiceType),
		),
	)
	return nil
}
