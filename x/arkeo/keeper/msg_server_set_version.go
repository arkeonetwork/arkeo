package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/errors"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) SetVersion(goCtx context.Context, msg *types.MsgSetVersion) (*types.MsgSetVersionResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	ctx.Logger().Info(
		"receive MsgSetVersion",
		"signer", msg.Creator,
		"version", msg.Version,
	)

	cacheCtx, commit := ctx.CacheContext()
	if err := k.SetVersionValidate(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed set version validation", "err", err)
		return nil, err
	}

	if err := k.SetVersionHandle(cacheCtx, msg); err != nil {
		ctx.Logger().Error("failed set version handle", "err", err)
		return nil, err
	}
	commit()

	return &types.MsgSetVersionResponse{}, nil
}

func (k msgServer) SetVersionValidate(ctx cosmos.Context, msg *types.MsgSetVersion) error {
	if k.FetchConfig(ctx, configs.HandlerSetVersion) > 0 {
		return errors.Wrapf(types.ErrDisabledHandler, "set version")
	}

	acct, _ := sdk.AccAddressFromBech32(msg.Creator)
	valAddr := cosmos.ValAddress(acct)
	currentVersion := k.GetVersionForAddress(ctx, valAddr)
	if currentVersion > msg.Version {
		return fmt.Errorf("cannot downgrade version: (%d/%d)", msg.Version, currentVersion)
	}
	if currentVersion == msg.Version {
		return fmt.Errorf("cannot set version to the same version: (%d/%d)", msg.Version, currentVersion)
	}

	return nil
}

func (k msgServer) SetVersionHandle(ctx cosmos.Context, msg *types.MsgSetVersion) error {
	acct, _ := sdk.AccAddressFromBech32(msg.Creator)
	valAddr := cosmos.ValAddress(acct)
	k.SetVersionForAddress(ctx, valAddr, msg.Version)
	return nil
}
