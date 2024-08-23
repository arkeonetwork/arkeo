package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptoCodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgBondProvider{}, "arkeo/BondProvider", nil)
	cdc.RegisterConcrete(&MsgModProvider{}, "arkeo/ModProvider", nil)
	cdc.RegisterConcrete(&MsgOpenContract{}, "arkeo/OpenContract", nil)
	cdc.RegisterConcrete(&MsgCloseContract{}, "arkeo/CloseContract", nil)
	cdc.RegisterConcrete(&MsgClaimContractIncome{}, "arkeo/ClaimContractIncome", nil)
	cdc.RegisterConcrete(&MsgSetVersion{}, "arkeo/SetVersion", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBondProvider{},
		&MsgModProvider{},
		&MsgOpenContract{},
		&MsgCloseContract{},
		&MsgClaimContractIncome{},
		&MsgSetVersion{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)

func init() {
	RegisterLegacyAminoCodec(Amino)
	cryptoCodec.RegisterCrypto(Amino)
	sdk.RegisterLegacyAminoCodec(Amino)
	Amino.Seal()
}
