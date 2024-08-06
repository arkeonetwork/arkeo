package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrAirdropEnded                = errors.Register(ModuleName, 1, "Airdrop has ended")
	ErrNoClaimableAmount           = errors.Register(ModuleName, 2, "No Claimable Arkeo")
	ErrInvalidSignature            = errors.Register(ModuleName, 3, "Invalid signature")
	ErrClaimRecordNotTransferrable = errors.Register(ModuleName, 4, "Claim record can not be transferred")
	ErrInvalidCreator              = errors.Register(ModuleName, 5, "Invalid Creator")
)
