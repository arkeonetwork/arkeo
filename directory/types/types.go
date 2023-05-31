package types

import "github.com/arkeonetwork/arkeo/common/cosmos"

type BondProviderEvent struct {
	Pubkey       string `mapstructure:"provider"`
	Service      string `mapstructure:"service"`
	Height       int64  `mapstructure:"height"`
	TxID         string `mapstructure:"hash"`
	BondRelative string `mapstructure:"bond_rel"`
	BondAbsolute string `mapstructure:"bond_abs"`
}

type ContractType string

var (
	ContractTypePayAsYouGo   ContractType = "PayAsYouGo"
	ContractTypeSubscription ContractType = "Subscription"
)

type AuthType string

var (
	AuthTypeStrict AuthType = "STRICT"
	AuthTypeOpen   AuthType = "OPEN"
)

type CloseContractEvent struct {
	ContractId     uint64 `mapstructure:"contract_id"`
	ProviderPubkey string `mapstructure:"provider"`
	Service        string `mapstructure:"service"`
	ClientPubkey   string `mapstructure:"client"`
	DelegatePubkey string `mapstructure:"delegate"`
	TxID           string `mapstructure:"hash"`
	Height         int64  `mapstructure:"height"`
	EventHeight    int64  `mapstructure:"eventHeight"`
}

// get the delegate pubkey falling back to client pubkey if undefined
func (c CloseContractEvent) GetDelegatePubkey() string {
	if c.DelegatePubkey != "" {
		return c.DelegatePubkey
	}
	return c.ClientPubkey
}

type ProviderStatus string

var (
	ProviderStatusOnline  ProviderStatus = "Online"
	ProviderStatusOffline ProviderStatus = "Offline"
)

type ModProviderEvent struct {
	Pubkey              string         `mapstructure:"provider"`
	Service             string         `mapstructure:"service"`
	Height              int64          `mapstructure:"height"`
	TxID                string         `mapstructure:"hash"`
	MetadataURI         string         `mapstructure:"metadata_uri"`
	MetadataNonce       uint64         `mapstructure:"metadata_nonce"`
	Status              ProviderStatus `mapstructure:"status"`
	MinContractDuration int64          `mapstructure:"min_contract_duration"`
	MaxContractDuration int64          `mapstructure:"max_contract_duration"`
	SettlementDuration  int64          `mapstructure:"settlement_duration"`
	SubscriptionRate    cosmos.Coins   `mapstructure:"subscription_rate"`
	PayAsYouGoRate      cosmos.Coins   `mapstructure:"pay_as_you_go_rate"`
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type ProviderSortKey string

var (
	ProviderSortKeyNone          ProviderSortKey = ""
	ProviderSortKeyAge           ProviderSortKey = "age"
	ProviderSortKeyContractCount ProviderSortKey = "contract_count"
	ProviderSortKeyAmountPaid    ProviderSortKey = "amount_paid"
)

type ProviderSearchParams struct {
	Pubkey                     string
	Service                    string
	SortKey                    ProviderSortKey
	MaxDistance                int64
	IsMaxDistanceSet           bool
	Coordinates                Coordinates
	MinValidatorPayments       int64
	IsMinValidatorPaymentsSet  bool
	MinProviderAge             int64
	IsMinProviderAgeSet        bool
	MinFreeRateLimit           int64
	IsMinFreeRateLimitSet      bool
	MinPaygoRateLimit          int64
	IsMinPaygoRateLimitSet     bool
	MinSubscribeRateLimit      int64
	IsMinSubscribeRateLimitSet bool
	MinOpenContracts           int64
	IsMinOpenContractsSet      bool
}

// swagger:model ArkeoStats
type ArkeoStats struct {
	ContractsOpen           int64 `db:"open_contracts"`
	ContractsTotal          int64 `db:"total_contracts"`
	ContractsMedianDuration int64 `db:"median_open_contract_length"`
	ContractsMedianRate     int64 `db:"median_open_contract_rate"`
	ProviderCount           int64 `db:"total_online_providers"`
	QueryCount              int64 `db:"total_queries"`
	TotalIncome             int64 `db:"total_paid"`
	// TODO: in the future we can add more complicated structure
	// ContractsMedianRatePayPer       int64
	// ContractsMedianRateSubscription int64
	// ServiceStats                      map[string]*ServiceStats
}

// swagger:model ServiceStats
type ServiceStats struct {
	Service            string
	ProviderCount      int64
	QueryCount         int64
	QueryCountLastDay  int64
	TotalIncome        int64
	TotalIncomeLastDay int64
}
