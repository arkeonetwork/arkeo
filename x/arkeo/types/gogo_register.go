package types

import (
	cosmosproto "github.com/cosmos/gogoproto/proto"
)

// Re-register message types into the cosmos/gogoproto registry so
// cdctypes.MsgTypeURL derives stable names even though protoc-gen-gocosmos
// registers against github.com/gogo/protobuf/proto.
func init() {
	cosmosproto.RegisterFile("arkeo/arkeo/tx.proto", fileDescriptor_a12700967a3e4015)
	cosmosproto.RegisterFile("arkeo/arkeo/query.proto", fileDescriptor_4b28dca1d1dd051d)
	cosmosproto.RegisterFile("arkeo/arkeo/genesis.proto", fileDescriptor_caae968dd754c6d4)
	cosmosproto.RegisterFile("arkeo/arkeo/keeper.proto", fileDescriptor_f833050061122841)
	cosmosproto.RegisterFile("arkeo/arkeo/misc.proto", fileDescriptor_64a3fa5463db3b34)
	cosmosproto.RegisterFile("arkeo/arkeo/params.proto", fileDescriptor_47c871f4fc73dfc5)
	cosmosproto.RegisterType((*MsgBondProvider)(nil), "arkeo.arkeo.MsgBondProvider")
	cosmosproto.RegisterType((*MsgBondProviderResponse)(nil), "arkeo.arkeo.MsgBondProviderResponse")
	cosmosproto.RegisterType((*MsgModProvider)(nil), "arkeo.arkeo.MsgModProvider")
	cosmosproto.RegisterType((*MsgModProviderResponse)(nil), "arkeo.arkeo.MsgModProviderResponse")
	cosmosproto.RegisterType((*MsgOpenContract)(nil), "arkeo.arkeo.MsgOpenContract")
	cosmosproto.RegisterType((*MsgOpenContractResponse)(nil), "arkeo.arkeo.MsgOpenContractResponse")
	cosmosproto.RegisterType((*MsgCloseContract)(nil), "arkeo.arkeo.MsgCloseContract")
	cosmosproto.RegisterType((*MsgCloseContractResponse)(nil), "arkeo.arkeo.MsgCloseContractResponse")
	cosmosproto.RegisterType((*MsgClaimContractIncome)(nil), "arkeo.arkeo.MsgClaimContractIncome")
	cosmosproto.RegisterType((*MsgClaimContractIncomeResponse)(nil), "arkeo.arkeo.MsgClaimContractIncomeResponse")
	cosmosproto.RegisterType((*MsgSetVersion)(nil), "arkeo.arkeo.MsgSetVersion")
	cosmosproto.RegisterType((*MsgSetVersionResponse)(nil), "arkeo.arkeo.MsgSetVersionResponse")
	cosmosproto.RegisterType((*MsgRegisterService)(nil), "arkeo.arkeo.MsgRegisterService")
	cosmosproto.RegisterType((*MsgRegisterServiceResponse)(nil), "arkeo.arkeo.MsgRegisterServiceResponse")
	cosmosproto.RegisterType((*MsgUpdateService)(nil), "arkeo.arkeo.MsgUpdateService")
	cosmosproto.RegisterType((*MsgUpdateServiceResponse)(nil), "arkeo.arkeo.MsgUpdateServiceResponse")
	cosmosproto.RegisterType((*MsgRemoveService)(nil), "arkeo.arkeo.MsgRemoveService")
	cosmosproto.RegisterType((*MsgRemoveServiceResponse)(nil), "arkeo.arkeo.MsgRemoveServiceResponse")
}
