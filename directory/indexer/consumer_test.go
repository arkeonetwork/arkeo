package indexer

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	abcitypes "github.com/tendermint/tendermint/abci/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

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

func TestParseEventToConcreteType(t *testing.T) {
	cosmos.GetConfig().SetBech32PrefixForAccount("tarkeo", "tarkeopub")
	input := `{"type":"arkeo.arkeo.EventOpenContract","attributes":[{"key":"YXV0aG9yaXphdGlvbg==","value":"IlNUUklDVCI=","index":true},{"key":"Y2xpZW50","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFncGpncDV2OHRqNmdkaDZnY3pxd3d3NWtzaDRnOHluYzh4cGpqc3Nhd3NuN2N4cXdtaG1qeTRkOGQ4Ig==","index":true},{"key":"Y29udHJhY3RfaWQ=","value":"IjIi","index":true},{"key":"ZGVsZWdhdGU=","value":"IiI=","index":true},{"key":"ZGVwb3NpdA==","value":"IjkwMCI=","index":true},{"key":"ZHVyYXRpb24=","value":"IjYwIg==","index":true},{"key":"aGVpZ2h0","value":"IjE0OTUi","index":true},{"key":"b3Blbl9jb3N0","value":"IjEwMDAwMDAwMCI=","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFmMHZtZ2h1YWtlZjR6eG5oNmh2Mmdld21xZ201dGRnOWY2dzNxeGpwdzQ5eG5zamYzNmY3ZjQwZXZlIg==","index":true},{"key":"cXVlcmllc19wZXJfbWludXRl","value":"IjEwIg==","index":true},{"key":"cmF0ZQ==","value":"eyJkZW5vbSI6InVhcmtlbyIsImFtb3VudCI6IjE1In0=","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true},{"key":"c2V0dGxlbWVudF9kdXJhdGlvbg==","value":"IjEwIg==","index":true},{"key":"dHlwZQ==","value":"IlBBWV9BU19ZT1VfR08i","index":true}]}`
	var event abcitypes.Event
	if err := json.Unmarshal([]byte(input), &event); err != nil {
		t.Error(err)
	}
	contractOpenEvent, err := parseEventToConcreteType[types.EventOpenContract](event)
	assert.Nil(t, err)
	assert.True(t, contractOpenEvent.ContractId == 2)
	assert.True(t, contractOpenEvent.Service == "mock")
}

func TestConvertContractSettlementEvent(t *testing.T) {
	cosmos.GetConfig().SetBech32PrefixForAccount("tarkeo", "tarkeopub")
	input := `{"type":"arkeo.arkeo.EventSettleContract","attributes":[{"key":"Y2xpZW50","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFncGpncDV2OHRqNmdkaDZnY3pxd3d3NWtzaDRnOHluYzh4cGpqc3Nhd3NuN2N4cXdtaG1qeTRkOGQ4Ig==","index":true},{"key":"Y29udHJhY3RfaWQ=","value":"IjIi","index":true},{"key":"ZGVsZWdhdGU=","value":"IiI=","index":true},{"key":"aGVpZ2h0","value":"IjE0OTUi","index":true},{"key":"bm9uY2U=","value":"IjAi","index":true},{"key":"cGFpZA==","value":"IjAi","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFmMHZtZ2h1YWtlZjR6eG5oNmh2Mmdld21xZ201dGRnOWY2dzNxeGpwdzQ5eG5zamYzNmY3ZjQwZXZlIg==","index":true},{"key":"cmVzZXJ2ZQ==","value":"IjAi","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true},{"key":"dHlwZQ==","value":"IlBBWV9BU19ZT1VfR08i","index":true}]}`
	var event abcitypes.Event
	if err := json.Unmarshal([]byte(input), &event); err != nil {
		t.Error(err)
	}
	contractSettlementEvent, err := parseEventToConcreteType[types.EventSettleContract](event)
	assert.Nil(t, err)
	assert.True(t, contractSettlementEvent.ContractId == 2)
}

func TestConvertEventModProvider(t *testing.T) {
	input := `{"type":"arkeo.arkeo.EventModProvider","attributes":[{"key":"Ym9uZA==","value":"IjIwMDAwMDAwMDAwIg==","index":true},{"key":"Y3JlYXRvcg==","value":"InRhcmtlbzE5MzU4ejI2andoM2U0cmQ2cHN4cWY4cTZmM3BlNmY4czd2MHgyYSI=","index":true},{"key":"bWF4X2NvbnRyYWN0X2R1cmF0aW9u","value":"IjEwMCI=","index":true},{"key":"bWV0YWRhdGFfbm9uY2U=","value":"IjEi","index":true},{"key":"bWV0YWRhdGFfdXJp","value":"Imh0dHA6Ly9sb2NhbGhvc3Q6MzYzNi9tZXRhZGF0YS5qc29uIg==","index":true},{"key":"bWluX2NvbnRyYWN0X2R1cmF0aW9u","value":"IjEwIg==","index":true},{"key":"cGF5X2FzX3lvdV9nb19yYXRl","value":"W3siZGVub20iOiJ1YXJrZW8iLCJhbW91bnQiOiIxNSJ9XQ==","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFmMHZtZ2h1YWtlZjR6eG5oNmh2Mmdld21xZ201dGRnOWY2dzNxeGpwdzQ5eG5zamYzNmY3ZjQwZXZlIg==","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true},{"key":"c2V0dGxlbWVudF9kdXJhdGlvbg==","value":"IjEwIg==","index":true},{"key":"c3RhdHVz","value":"Ik9OTElORSI=","index":true},{"key":"c3Vic2NyaXB0aW9uX3JhdGU=","value":"W3siZGVub20iOiJ1YXJrZW8iLCJhbW91bnQiOiIxMCJ9XQ==","index":true}]}`
	cosmos.GetConfig().SetBech32PrefixForAccount("tarkeo", "tarkeopub")
	var event abcitypes.Event
	if err := json.Unmarshal([]byte(input), &event); err != nil {
		t.Error(err)
	}
	eventModProvider, err := parseEventToConcreteType[types.EventModProvider](event)
	assert.Nil(t, err)
	assert.True(t, eventModProvider.Service == "mock")
}
