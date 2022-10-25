package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/mercury module sentinel errors
var (
	ErrProviderBadSigner     = sdkerrors.Register(ModuleName, 2, "unauthorized: bad provider pubkey and signer")
	ErrProviderAlreadyExists = sdkerrors.Register(ModuleName, 3, "provider already exists")
	ErrInsufficientFunds     = sdkerrors.Register(ModuleName, 4, "insufficient funds")
	ErrInvalidBond           = sdkerrors.Register(ModuleName, 5, "invalid bond")
)
