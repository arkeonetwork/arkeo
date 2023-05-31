package indexer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	abcitypes "github.com/tendermint/tendermint/abci/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	arkeotypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestEventParsing(t *testing.T) {
	cosmos.GetConfig().SetBech32PrefixForAccount("tarkeo", "tarkeopub")
	inputs := []struct {
		Name    string
		Payload string
		Checker func(t *testing.T, result any)
	}{
		{
			Name:    "EventOpenContract",
			Payload: `{"type":"arkeo.arkeo.EventOpenContract","attributes":[{"key":"YXV0aG9yaXphdGlvbg==","value":"IlNUUklDVCI=","index":true},{"key":"Y2xpZW50","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFncGpncDV2OHRqNmdkaDZnY3pxd3d3NWtzaDRnOHluYzh4cGpqc3Nhd3NuN2N4cXdtaG1qeTRkOGQ4Ig==","index":true},{"key":"Y29udHJhY3RfaWQ=","value":"IjIi","index":true},{"key":"ZGVsZWdhdGU=","value":"IiI=","index":true},{"key":"ZGVwb3NpdA==","value":"IjkwMCI=","index":true},{"key":"ZHVyYXRpb24=","value":"IjYwIg==","index":true},{"key":"aGVpZ2h0","value":"IjE0OTUi","index":true},{"key":"b3Blbl9jb3N0","value":"IjEwMDAwMDAwMCI=","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFmMHZtZ2h1YWtlZjR6eG5oNmh2Mmdld21xZ201dGRnOWY2dzNxeGpwdzQ5eG5zamYzNmY3ZjQwZXZlIg==","index":true},{"key":"cXVlcmllc19wZXJfbWludXRl","value":"IjEwIg==","index":true},{"key":"cmF0ZQ==","value":"eyJkZW5vbSI6InVhcmtlbyIsImFtb3VudCI6IjE1In0=","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true},{"key":"c2V0dGxlbWVudF9kdXJhdGlvbg==","value":"IjEwIg==","index":true},{"key":"dHlwZQ==","value":"IlBBWV9BU19ZT1VfR08i","index":true}]}`,
			Checker: func(t *testing.T, result any) {
				assert.IsType(t, arkeotypes.EventOpenContract{}, result)
				e, ok := result.(arkeotypes.EventOpenContract)
				assert.True(t, ok)
				assert.Equal(t, "mock", e.Service)
				assert.Equal(t, "tarkeopub1addwnpepqf0vmghuakef4zxnh6hv2gewmqgm5tdg9f6w3qxjpw49xnsjf36f7f40eve", e.Provider.String())
				assert.Equal(t, "tarkeopub1addwnpepqgpjgp5v8tj6gdh6gczqwww5ksh4g8ync8xpjjssawsn7cxqwmhmjy4d8d8", e.Client.String())
				assert.Equal(t, uint64(2), e.ContractId)
				assert.Nil(t, e.Delegate)
				assert.Equal(t, int64(1495), e.Height)
				assert.Equal(t, int64(60), e.Duration)
				assert.Equal(t, int64(100000000), e.OpenCost)
				assert.Equal(t, int64(10), e.QueriesPerMinute)
				assert.Equal(t, arkeotypes.ContractAuthorization_STRICT, e.Authorization)
				assert.Equal(t, arkeotypes.ContractType_PAY_AS_YOU_GO, e.Type)
				assert.NotNil(t, e.Rate)
				assert.Equal(t, "uarkeo", e.Rate.Denom)
				assert.Equal(t, int64(15), e.Rate.Amount.Int64())
			},
		},
		{
			Name:    "EventSettleContract",
			Payload: `{"type":"arkeo.arkeo.EventSettleContract","attributes":[{"key":"Y2xpZW50","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFncGpncDV2OHRqNmdkaDZnY3pxd3d3NWtzaDRnOHluYzh4cGpqc3Nhd3NuN2N4cXdtaG1qeTRkOGQ4Ig==","index":true},{"key":"Y29udHJhY3RfaWQ=","value":"IjIi","index":true},{"key":"ZGVsZWdhdGU=","value":"IiI=","index":true},{"key":"aGVpZ2h0","value":"IjE0OTUi","index":true},{"key":"bm9uY2U=","value":"IjAi","index":true},{"key":"cGFpZA==","value":"IjAi","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFmMHZtZ2h1YWtlZjR6eG5oNmh2Mmdld21xZ201dGRnOWY2dzNxeGpwdzQ5eG5zamYzNmY3ZjQwZXZlIg==","index":true},{"key":"cmVzZXJ2ZQ==","value":"IjAi","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true},{"key":"dHlwZQ==","value":"IlBBWV9BU19ZT1VfR08i","index":true}]}`,
			Checker: func(t *testing.T, result any) {
				assert.IsType(t, arkeotypes.EventSettleContract{}, result)
				e, ok := result.(arkeotypes.EventSettleContract)
				assert.True(t, ok)
				assert.Equal(t, "tarkeopub1addwnpepqf0vmghuakef4zxnh6hv2gewmqgm5tdg9f6w3qxjpw49xnsjf36f7f40eve", e.Provider.String())
				assert.Equal(t, uint64(2), e.ContractId)
				assert.Equal(t, "mock", e.Service)
				assert.Equal(t, "tarkeopub1addwnpepqgpjgp5v8tj6gdh6gczqwww5ksh4g8ync8xpjjssawsn7cxqwmhmjy4d8d8", e.Client.String())
				assert.Equal(t, arkeotypes.ContractType_PAY_AS_YOU_GO, e.Type)
				assert.Nil(t, e.Delegate)
				assert.Zero(t, e.Nonce)
				assert.Equal(t, int64(1495), e.Height)
				assert.Equal(t, int64(0), e.Paid.Int64())
				assert.Equal(t, int64(0), e.Reserve.Int64())
			},
		},
		{
			Name:    "EventModProvider",
			Payload: `{"type":"arkeo.arkeo.EventModProvider","attributes":[{"key":"Ym9uZA==","value":"IjIwMDAwMDAwMDAwIg==","index":true},{"key":"Y3JlYXRvcg==","value":"InRhcmtlbzE5MzU4ejI2andoM2U0cmQ2cHN4cWY4cTZmM3BlNmY4czd2MHgyYSI=","index":true},{"key":"bWF4X2NvbnRyYWN0X2R1cmF0aW9u","value":"IjEwMCI=","index":true},{"key":"bWV0YWRhdGFfbm9uY2U=","value":"IjEi","index":true},{"key":"bWV0YWRhdGFfdXJp","value":"Imh0dHA6Ly9sb2NhbGhvc3Q6MzYzNi9tZXRhZGF0YS5qc29uIg==","index":true},{"key":"bWluX2NvbnRyYWN0X2R1cmF0aW9u","value":"IjEwIg==","index":true},{"key":"cGF5X2FzX3lvdV9nb19yYXRl","value":"W3siZGVub20iOiJ1YXJrZW8iLCJhbW91bnQiOiIxNSJ9XQ==","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFmMHZtZ2h1YWtlZjR6eG5oNmh2Mmdld21xZ201dGRnOWY2dzNxeGpwdzQ5eG5zamYzNmY3ZjQwZXZlIg==","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true},{"key":"c2V0dGxlbWVudF9kdXJhdGlvbg==","value":"IjEwIg==","index":true},{"key":"c3RhdHVz","value":"Ik9OTElORSI=","index":true},{"key":"c3Vic2NyaXB0aW9uX3JhdGU=","value":"W3siZGVub20iOiJ1YXJrZW8iLCJhbW91bnQiOiIxMCJ9XQ==","index":true}]}`,
			Checker: func(t *testing.T, result any) {
				assert.IsType(t, arkeotypes.EventModProvider{}, result)
				e, ok := result.(arkeotypes.EventModProvider)
				assert.True(t, ok)
				assert.Equal(t, "mock", e.Service)
				assert.Equal(t, "tarkeopub1addwnpepqf0vmghuakef4zxnh6hv2gewmqgm5tdg9f6w3qxjpw49xnsjf36f7f40eve", e.Provider.String())
				assert.Equal(t, "http://localhost:3636/metadata.json", e.MetadataUri)
				assert.Equal(t, uint64(1), e.MetadataNonce)
				assert.Equal(t, arkeotypes.ProviderStatus_ONLINE, e.Status)
				assert.Equal(t, int64(10), e.MinContractDuration)
				assert.Equal(t, int64(100), e.MaxContractDuration)
				assert.Len(t, e.SubscriptionRate, 1)
				assert.Equal(t, "uarkeo", e.SubscriptionRate[0].Denom)
				assert.Equal(t, int64(10), e.SubscriptionRate[0].Amount.Int64())
				assert.Len(t, e.PayAsYouGoRate, 1)
				assert.Equal(t, "uarkeo", e.PayAsYouGoRate[0].Denom)
				assert.Equal(t, int64(15), e.PayAsYouGoRate[0].Amount.Int64())
				assert.Equal(t, int64(20000000000), e.Bond.Int64())
				assert.Equal(t, int64(10), e.SettlementDuration)
			},
		},
		{
			Name:    "EventBondProvider",
			Payload: `{"type":"arkeo.arkeo.EventBondProvider","attributes":[{"key":"Ym9uZF9hYnM=","value":"IjEwMDAwMDAwMDAi","index":true},{"key":"Ym9uZF9yZWw=","value":"IjEwMDAwMDAwMDAi","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFncGpncDV2OHRqNmdkaDZnY3pxd3d3NWtzaDRnOHluYzh4cGpqc3Nhd3NuN2N4cXdtaG1qeTRkOGQ4Ig==","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true}]}`,
			Checker: func(t *testing.T, result any) {
				assert.IsType(t, arkeotypes.EventBondProvider{}, result)
				e, ok := result.(arkeotypes.EventBondProvider)
				assert.True(t, ok)
				assert.Equal(t, "mock", e.Service)
				assert.Equal(t, "tarkeopub1addwnpepqgpjgp5v8tj6gdh6gczqwww5ksh4g8ync8xpjjssawsn7cxqwmhmjy4d8d8", e.Provider.String())
				assert.NotNil(t, e.BondRel)
				assert.NotNil(t, e.BondAbs)
				assert.Equal(t, int64(1000000000), e.BondRel.Int64())
				assert.Equal(t, int64(1000000000), e.BondAbs.Int64())
			},
		},
	}
	for _, c := range inputs {
		var event abcitypes.Event
		if err := json.Unmarshal([]byte(c.Payload), &event); err != nil {
			t.Error(err)
		}
		var result any
		var err error
		switch event.Type {
		case arkeotypes.EventTypeOpenContract:
			result, err = parseEventToConcreteType[arkeotypes.EventOpenContract](event)
		case arkeotypes.EventTypeSettleContract:
			result, err = parseEventToConcreteType[arkeotypes.EventSettleContract](event)
		case arkeotypes.EventTypeModProvider:
			result, err = parseEventToConcreteType[arkeotypes.EventModProvider](event)
		case arkeotypes.EventTypeBondProvider:
			result, err = parseEventToConcreteType[arkeotypes.EventBondProvider](event)
		}
		assert.Nil(t, err)
		c.Checker(t, result)
	}
}

func TestConvertEventToMap(t *testing.T) {
	cosmos.GetConfig().SetBech32PrefixForAccount("tarkeo", "tarkeopub")
	input := `{"type":"arkeo.arkeo.EventOpenContract","attributes":[{"key":"YXV0aG9yaXphdGlvbg==","value":"IlNUUklDVCI=","index":true},{"key":"Y2xpZW50","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFncGpncDV2OHRqNmdkaDZnY3pxd3d3NWtzaDRnOHluYzh4cGpqc3Nhd3NuN2N4cXdtaG1qeTRkOGQ4Ig==","index":true},{"key":"Y29udHJhY3RfaWQ=","value":"IjIi","index":true},{"key":"ZGVsZWdhdGU=","value":"IiI=","index":true},{"key":"ZGVwb3NpdA==","value":"IjkwMCI=","index":true},{"key":"ZHVyYXRpb24=","value":"IjYwIg==","index":true},{"key":"aGVpZ2h0","value":"IjE0OTUi","index":true},{"key":"b3Blbl9jb3N0","value":"IjEwMDAwMDAwMCI=","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFmMHZtZ2h1YWtlZjR6eG5oNmh2Mmdld21xZ201dGRnOWY2dzNxeGpwdzQ5eG5zamYzNmY3ZjQwZXZlIg==","index":true},{"key":"cXVlcmllc19wZXJfbWludXRl","value":"IjEwIg==","index":true},{"key":"cmF0ZQ==","value":"eyJkZW5vbSI6InVhcmtlbyIsImFtb3VudCI6IjE1In0=","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true},{"key":"c2V0dGxlbWVudF9kdXJhdGlvbg==","value":"IjEwIg==","index":true},{"key":"dHlwZQ==","value":"IlBBWV9BU19ZT1VfR08i","index":true}]}`
	var event abcitypes.Event
	if err := json.Unmarshal([]byte(input), &event); err != nil {
		t.Error(err)
	}
	result, err := convertEventToMap(event)
	assert.Nil(t, err)
	rate, ok := result["rate"]
	assert.True(t, ok)
	assert.IsType(t, map[string]any{}, rate)
}
