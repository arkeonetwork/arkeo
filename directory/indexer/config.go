package indexer

import "github.com/arkeonetwork/arkeo/directory/db"

// ServiceParams hold all necessary parameters for indexer app to run
type ServiceParams struct {
	ArkeoApi            string      `mapstructure:"arkeo_api" json:"arkeo_api"`
	TendermintApi       string      `mapstructure:"tendermint_api" json:"tendermint_api"`
	TendermintWs        string      `mapstructure:"tendermint_ws" json:"tendermint_ws"`
	ChainID             string      `mapstructure:"chain_id" json:"chain_id"`
	Bech32PrefixAccAddr string      `mapstructure:"bech32_pref_acc_addr" json:"bech32_pref_acc_addr"`
	Bech32PrefixAccPub  string      `mapstructure:"bech32_pref_acc_pub" json:"bech32_pref_acc_pub"`
	IndexerID           int64       `json:"-"`
	DB                  db.DBConfig `mapstructure:"db" json:"db"`
}

//
//import (
//	"github.com/cosmos/cosmos-sdk/client"
//	"github.com/cosmos/cosmos-sdk/codec"
//	"github.com/cosmos/cosmos-sdk/codec/types"
//	stdtypes "github.com/cosmos/cosmos-sdk/std"
//	"github.com/cosmos/cosmos-sdk/x/auth/tx"
//	authztypes "github.com/cosmos/cosmos-sdk/x/authz"
//	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
//	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
//	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
//	// ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
//	// ibccoretypes "github.com/cosmos/ibc-go/v3/modules/core/types"
//)
//
//type encodingConfig struct {
//	InterfaceRegistry types.InterfaceRegistry
//	Marshaler         codec.Codec
//	TxConfig          client.TxConfig
//	Amino             *codec.LegacyAmino
//}
//
//// NewEncoding registers all base protobuf types by default as well as any custom types passed in
//func NewEncoding(registerInterfaces ...func(r types.InterfaceRegistry)) *encodingConfig {
//	registry := types.NewInterfaceRegistry()
//
//	// register base protobuf types
//	authztypes.RegisterInterfaces(registry)
//	banktypes.RegisterInterfaces(registry)
//	distributiontypes.RegisterInterfaces(registry)
//	// ibccoretypes.RegisterInterfaces(registry)
//	// ibctransfertypes.RegisterInterfaces(registry)
//	stakingtypes.RegisterInterfaces(registry)
//	stdtypes.RegisterInterfaces(registry)
//
//	// register custom protobuf types
//	for _, r := range registerInterfaces {
//		r(registry)
//	}
//
//	marshaler := codec.NewProtoCodec(registry)
//
//	return &encodingConfig{
//		InterfaceRegistry: registry,
//		Marshaler:         marshaler,
//		TxConfig:          tx.NewTxConfig(marshaler, tx.DefaultSignModes),
//		Amino:             codec.NewLegacyAmino(),
//	}
//}
