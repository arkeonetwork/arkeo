package utils

import (
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

// create a codec used only for testing
func MakeTestCodec() *codec.LegacyAmino {
	cdc := codec.NewLegacyAmino()
	banktypes.RegisterLegacyAminoCodec(cdc)
	authtypes.RegisterLegacyAminoCodec(cdc)
	types.RegisterCodec(cdc)
	cosmos.RegisterCodec(cdc)
	// codec.RegisterCrypto(cdc)
	return cdc
}

// GetRandomArkeoAddress returns a random arkeo address for testing
func GetRandomArkeoAddress() sdk.AccAddress {
	return sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
}
