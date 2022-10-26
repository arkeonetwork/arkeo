package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/mercury module sentinel errors
var (
	ErrProviderBadSigner                     = sdkerrors.Register(ModuleName, 2, "unauthorized: bad provider pubkey and signer")
	ErrProviderAlreadyExists                 = sdkerrors.Register(ModuleName, 3, "provider already exists")
	ErrInsufficientFunds                     = sdkerrors.Register(ModuleName, 4, "insufficient funds")
	ErrInvalidBond                           = sdkerrors.Register(ModuleName, 5, "invalid bond")
	ErrInvalidModProviderMetdataURI          = sdkerrors.Register(ModuleName, 6, "invalid mod provider metadata uri")
	ErrInvalidModProviderMaxContractDuration = sdkerrors.Register(ModuleName, 7, "invalid mod provider max contract duration")
	ErrInvalidModProviderMinContractDuration = sdkerrors.Register(ModuleName, 8, "invalid mod provider min contract duration")
	ErrInvalidModProviderStatus              = sdkerrors.Register(ModuleName, 9, "invalid mod provider bad provider status")
	ErrDisabledHandler                       = sdkerrors.Register(ModuleName, 10, "disabled handler")
)
