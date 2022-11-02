package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

func RegisterCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgBondProvider{}, "mercury/BondProvider", nil)
	cdc.RegisterConcrete(&MsgModProvider{}, "mercury/ModProvider", nil)
	cdc.RegisterConcrete(&MsgOpenContract{}, "mercury/OpenContract", nil)
	cdc.RegisterConcrete(&MsgCloseContract{}, "mercury/CloseContract", nil)
	cdc.RegisterConcrete(&MsgClaimContractIncome{}, "mercury/ClaimContractIncome", nil)
	// this line is used by starport scaffolding # 2
}

func RegisterInterfaces(registry cdctypes.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgBondProvider{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgModProvider{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgOpenContract{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCloseContract{},
	)
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgClaimContractIncome{},
	)
	// this line is used by starport scaffolding # 3

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	Amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewProtoCodec(cdctypes.NewInterfaceRegistry())
)
