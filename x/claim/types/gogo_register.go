package types

import (
	cosmosproto "github.com/cosmos/gogoproto/proto"
)

// Mirror registrations into the cosmos/gogoproto registry so MsgTypeURL
// resolves the correct type URLs when using the SDK interface registry.
func init() {
	cosmosproto.RegisterFile("arkeo/claim/tx.proto", fileDescriptor_6a4ddac60cb43154)
	cosmosproto.RegisterFile("arkeo/claim/query.proto", fileDescriptor_c9bba750e7be3a92)
	cosmosproto.RegisterFile("arkeo/claim/genesis.proto", fileDescriptor_087fbceea1b57c13)
	cosmosproto.RegisterFile("arkeo/claim/params.proto", fileDescriptor_2bdbd8eba92221b6)
	cosmosproto.RegisterFile("arkeo/claim/claim_record.proto", fileDescriptor_db5386e8ec5cd28f)
	cosmosproto.RegisterType((*MsgClaimEth)(nil), "arkeo.claim.MsgClaimEth")
	cosmosproto.RegisterType((*MsgClaimEthResponse)(nil), "arkeo.claim.MsgClaimEthResponse")
	cosmosproto.RegisterType((*MsgClaimArkeo)(nil), "arkeo.claim.MsgClaimArkeo")
	cosmosproto.RegisterType((*MsgClaimArkeoResponse)(nil), "arkeo.claim.MsgClaimArkeoResponse")
	cosmosproto.RegisterType((*MsgTransferClaim)(nil), "arkeo.claim.MsgTransferClaim")
	cosmosproto.RegisterType((*MsgTransferClaimResponse)(nil), "arkeo.claim.MsgTransferClaimResponse")
	cosmosproto.RegisterType((*MsgAddClaim)(nil), "arkeo.claim.MsgAddClaim")
	cosmosproto.RegisterType((*MsgAddClaimResponse)(nil), "arkeo.claim.MsgAddClaimResponse")
	cosmosproto.RegisterType((*MsgClaimThorchain)(nil), "arkeo.claim.MsgClaimThorchain")
	cosmosproto.RegisterType((*MsgClaimThorchainResponse)(nil), "arkeo.claim.MsgClaimThorchainResponse")
}
