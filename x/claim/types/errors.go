package types

import (
	"cosmossdk.io/errors"
)

var (
	ErrAirdropEnded                = errors.Register(ModuleName, 2, "Airdrop has ended")
	ErrNoClaimableAmount           = errors.Register(ModuleName, 3, "No Claimable Arkeo")
	ErrInvalidSignature            = errors.Register(ModuleName, 4, "Invalid signature")
	ErrClaimRecordNotTransferrable = errors.Register(ModuleName, 5, "Claim record can not be transferred")
)
