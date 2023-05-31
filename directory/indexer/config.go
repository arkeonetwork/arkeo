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
