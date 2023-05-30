package indexer

import (
	"encoding/json"
	"fmt"
	"testing"

	abcitypes "github.com/tendermint/tendermint/abci/types"
)

func TestParseEventToEventOpenContract(t *testing.T) {
	input := `{"type":"arkeo.arkeo.EventOpenContract","attributes":[{"key":"YXV0aG9yaXphdGlvbg==","value":"IlNUUklDVCI=","index":true},{"key":"Y2xpZW50","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFncGpncDV2OHRqNmdkaDZnY3pxd3d3NWtzaDRnOHluYzh4cGpqc3Nhd3NuN2N4cXdtaG1qeTRkOGQ4Ig==","index":true},{"key":"Y29udHJhY3RfaWQ=","value":"IjIi","index":true},{"key":"ZGVsZWdhdGU=","value":"IiI=","index":true},{"key":"ZGVwb3NpdA==","value":"IjkwMCI=","index":true},{"key":"ZHVyYXRpb24=","value":"IjYwIg==","index":true},{"key":"aGVpZ2h0","value":"IjE0OTUi","index":true},{"key":"b3Blbl9jb3N0","value":"IjEwMDAwMDAwMCI=","index":true},{"key":"cHJvdmlkZXI=","value":"InRhcmtlb3B1YjFhZGR3bnBlcHFmMHZtZ2h1YWtlZjR6eG5oNmh2Mmdld21xZ201dGRnOWY2dzNxeGpwdzQ5eG5zamYzNmY3ZjQwZXZlIg==","index":true},{"key":"cXVlcmllc19wZXJfbWludXRl","value":"IjEwIg==","index":true},{"key":"cmF0ZQ==","value":"eyJkZW5vbSI6InVhcmtlbyIsImFtb3VudCI6IjE1In0=","index":true},{"key":"c2VydmljZQ==","value":"Im1vY2si","index":true},{"key":"c2V0dGxlbWVudF9kdXJhdGlvbg==","value":"IjEwIg==","index":true},{"key":"dHlwZQ==","value":"IlBBWV9BU19ZT1VfR08i","index":true}]}`
	var event abcitypes.Event
	if err := json.Unmarshal([]byte(input), &event); err != nil {
		t.Error(err)
	}
	result, err := parseEventToEventOpenContract(event)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%q", result)
	//result := make(map[string]interface{})
	//for _, att := range event.Attributes {
	//	if att.Value[0] == byte('{') {
	//		nest := make(map[string]interface{})
	//		if err := json.Unmarshal(att.Value, &nest); err != nil {
	//			panic(err)
	//		}
	//		result[string(att.Key)] = nest
	//	}
	//	result[string(att.Key)] = string(att.Value)
	//}
	//fmt.Printf("%+v\n", event)
	//var eventOpenContract types.EventOpenContract
	//decode, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
	//	Result:  &eventOpenContract,
	//	TagName: "json",
	//})
	//if err != nil {
	//	panic(err)
	//}
	//if err := decode.Decode(result); err != nil {
	//	panic(err)
	//}

}
