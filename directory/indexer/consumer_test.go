package indexer

import (
	"encoding/json"
	"fmt"
	"testing"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/assert"

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
			Payload: `{ "type": "arkeo.arkeo.EventOpenContract", "attributes": [ { "key": "authorization", "value": "\"STRICT\"", "index": true }, { "key": "client", "value": "\"tarkeopub1addwnpepqgpjgp5v8tj6gdh6gczqwww5ksh4g8ync8xpjjssawsn7cxqwmhmjy4d8d8\"", "index": true }, { "key": "contract_id", "value": "\"2\"", "index": true }, { "key": "delegate", "value": "\"\"", "index": true }, { "key": "deposit", "value": "\"900\"", "index": true }, { "key": "duration", "value": "\"60\"", "index": true }, { "key": "height", "value": "\"1495\"", "index": true }, { "key": "open_cost", "value": "\"100000000\"", "index": true }, { "key": "provider", "value": "\"tarkeopub1addwnpepqf0vmghuakef4zxnh6hv2gewmqgm5tdg9f6w3qxjpw49xnsjf36f7f40eve\"", "index": true }, { "key": "queries_per_minute", "value": "\"10\"", "index": true }, { "key": "rate", "value": "{\"denom\":\"uarkeo\",\"amount\":\"15\"}", "index": true }, { "key": "service", "value": "\"mock\"", "index": true }, { "key": "settlement_duration", "value": "\"10\"", "index": true }, { "key": "type", "value": "\"PAY_AS_YOU_GO\"", "index": true } ] }`,
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
			Payload: `{ "type": "arkeo.arkeo.EventSettleContract", "attributes": [ { "key": "client", "value": "\"tarkeopub1addwnpepqgpjgp5v8tj6gdh6gczqwww5ksh4g8ync8xpjjssawsn7cxqwmhmjy4d8d8\"", "index": true }, { "key": "contract_id", "value": "\"2\"", "index": true }, { "key": "delegate", "value": "\"\"", "index": true }, { "key": "height", "value": "\"1495\"", "index": true }, { "key": "nonce", "value": "\"0\"", "index": true }, { "key": "paid", "value": "\"0\"", "index": true }, { "key": "provider", "value": "\"tarkeopub1addwnpepqf0vmghuakef4zxnh6hv2gewmqgm5tdg9f6w3qxjpw49xnsjf36f7f40eve\"", "index": true }, { "key": "reserve", "value": "\"0\"", "index": true }, { "key": "service", "value": "\"mock\"", "index": true }, { "key": "type", "value": "\"PAY_AS_YOU_GO\"", "index": true } ] }`,
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
			Payload: `{ "type": "arkeo.arkeo.EventModProvider", "attributes": [ { "key": "bond", "value": "\"20000000000\"", "index": true }, { "key": "creator", "value": "\"tarkeo19358z26jwh3e4rd6psxqf8q6f3pe6f8s7v0x2a\"", "index": true }, { "key": "max_contract_duration", "value": "\"100\"", "index": true }, { "key": "metadata_nonce", "value": "\"1\"", "index": true }, { "key": "metadata_uri", "value": "\"http://localhost:3636/metadata.json\"", "index": true }, { "key": "min_contract_duration", "value": "\"10\"", "index": true }, { "key": "pay_as_you_go_rate", "value": "[{\"denom\":\"uarkeo\",\"amount\":\"15\"}]", "index": true }, { "key": "provider", "value": "\"tarkeopub1addwnpepqf0vmghuakef4zxnh6hv2gewmqgm5tdg9f6w3qxjpw49xnsjf36f7f40eve\"", "index": true }, { "key": "service", "value": "\"mock\"", "index": true }, { "key": "settlement_duration", "value": "\"10\"", "index": true }, { "key": "status", "value": "\"ONLINE\"", "index": true }, { "key": "subscription_rate", "value": "[{\"denom\":\"uarkeo\",\"amount\":\"10\"}]", "index": true } ] }`,
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
			Payload: `{"type": "arkeo.arkeo.EventBondProvider", "attributes": [ { "key": "bond_abs", "value": "\"1000000000\"", "index": true }, { "key": "bond_rel", "value": "\"1000000000\"", "index": true }, { "key": "provider", "value": "\"tarkeopub1addwnpepqgpjgp5v8tj6gdh6gczqwww5ksh4g8ync8xpjjssawsn7cxqwmhmjy4d8d8\"", "index": true }, { "key": "service", "value": "\"mock\"", "index": true } ] }`,
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

func TestConvertEventToMapString(t *testing.T) {
	cosmos.GetConfig().SetBech32PrefixForAccount("tarkeo", "tarkeopub")
	input := `{
  "type": "arkeo.arkeo.EventOpenContract",
  "attributes": [
    {
      "key": "authorization",
      "value": "STRICT",
      "index": true
    },
    {
      "key": "client",
      "value": "tarkeopub1addwnpepqgpjgp5v8tj6gdh6gczqwww5ksh4g8ync8xpjjssawsn7cxqwmhmjy4d8d8",
      "index": true
    },
    {
      "key": "contract_id",
      "value": "2",
      "index": true
    },
    {
      "key": "delegate",
      "value": "",
      "index": true
    },
    {
      "key": "deposit",
      "value": "900",
      "index": true
    },
    {
      "key": "duration",
      "value": "60",
      "index": true
    },
    {
      "key": "height",
      "value": "1495",
      "index": true
    },
    {
      "key": "open_cost",
      "value": "100000000",
      "index": true
    },
    {
      "key": "provider",
      "value": "tarkeopub1addwnpepqf0vmghuakef4zxnh6hv2gewmqgm5tdg9f6w3qxjpw49xnsjf36f7f40eve",
      "index": true
    },
    {
      "key": "queries_per_minute",
      "value": "10",
      "index": true
    },
    {
      "key": "rate",
      "value": "{\"denom\":\"uarkeo\",\"amount\":\"15\"}",
      "index": true
    },
    {
      "key": "service",
      "value": "mock",
      "index": true
    },
    {
      "key": "settlement_duration",
      "value": "10",
      "index": true
    },
    {
      "key": "type",
      "value": "PAY_AS_YOU_GO",
      "index": true
    }
  ]
}
`
	var event abcitypes.Event
	if err := json.Unmarshal([]byte(input), &event); err != nil {
		t.Error(err)
	}
	result, err := convertEventToMap(event)
	assert.Nil(t, err)
	rate, ok := result["rate"]
	assert.True(t, ok)
	fmt.Println(rate)
}
