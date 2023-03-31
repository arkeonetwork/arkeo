package sentinel

import (
	"fmt"
	"testing"
	"time"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
)

var testConfig = conf.Configuration{
	Moniker:                   "Testy McTestface",
	Website:                   "testing.com",
	Description:               "the best testnet ever",
	Location:                  "100,100",
	Port:                      "3636",
	ProxyHost:                 "localhost:3637",
	SourceChain:               "localhost", // this should point to arkeo rpc endpoints, but we can ignore for testing
	EventStreamHost:           "localhost",
	ProviderPubKey:            types.GetRandomPubKey(),
	FreeTierRateLimit:         10,
	FreeTierRateLimitDuration: time.Second,
	SubTierRateLimit:          10,
	SubTierRateLimitDuration:  time.Second,
	AsGoTierRateLimit:         10,
	AsGoTierRateLimitDuration: time.Second,
	ClaimStoreLocation:        "",
	GaiaRpcArchiveHost:        "gaia-host",
}

func TestHandleOpenContractEvent(t *testing.T) {
	proxy := NewProxy(testConfig)
	inputContract := types.Contract{
		Provider:           testConfig.ProviderPubKey,
		Service:            common.BTCService,
		Client:             types.GetRandomPubKey(),
		Delegate:           common.EmptyPubKey,
		Type:               types.ContractType_PAY_AS_YOU_GO,
		Height:             100,
		Duration:           100,
		Rate:               cosmos.NewInt64Coin("uarkeo", 1),
		Deposit:            sdk.NewInt(100),
		Nonce:              0,
		Id:                 1,
		SettlementDuration: 10,
	}
	openCost := int64(100)
	events := sdk.Events{
		types.NewOpenContractEvent(openCost, &inputContract),
	}
	proxy.handleOpenContractEvent(convertEventsToResultEvent(events))

	// confirm that our memstore has the contract and its active
	outputContract, err := proxy.MemStore.Get(inputContract.Key())
	require.NoError(t, err)
	require.Equal(t, inputContract, outputContract)
	_, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Service, inputContract.Client)
	require.NoError(t, err)

	// confirm that a contract for a different provider doesn't get stored.
	inputContract.Provider = types.GetRandomPubKey()
	inputContract.Id = 2
	events = sdk.Events{
		types.NewOpenContractEvent(openCost, &inputContract),
	}
	proxy.handleOpenContractEvent(convertEventsToResultEvent(events))
	_, err = proxy.MemStore.Get(inputContract.Key())
	require.Error(t, err)

	// confirm that we return the correct active contract when multiple contracts are present.
	proxy.MemStore.SetHeight(201) // contract with id 1 should now be expired

	_, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Service, inputContract.Client)
	require.Error(t, err) // contract not found since its expired.

	inputContract.Provider = testConfig.ProviderPubKey
	inputContract.Id = 3
	inputContract.Height = 200
	events = sdk.Events{
		types.NewOpenContractEvent(openCost, &inputContract),
	}
	proxy.handleOpenContractEvent(convertEventsToResultEvent(events))

	outputContract, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Service, inputContract.Client)
	require.NoError(t, err)
	require.Equal(t, inputContract, outputContract)
}

func TestHandleCloseContractEvent(t *testing.T) {
	proxy := NewProxy(testConfig)
	inputContract := types.Contract{
		Provider:           testConfig.ProviderPubKey,
		Service:            common.BTCService,
		Client:             types.GetRandomPubKey(),
		Delegate:           common.EmptyPubKey,
		Type:               types.ContractType_PAY_AS_YOU_GO,
		Height:             100,
		Duration:           100,
		Rate:               cosmos.NewInt64Coin("uarkeo", 1),
		Deposit:            sdk.NewInt(100),
		Nonce:              0,
		Id:                 1,
		SettlementDuration: 10,
	}
	openCost := int64(100)
	events := sdk.Events{
		types.NewOpenContractEvent(openCost, &inputContract),
	}
	proxy.handleOpenContractEvent(convertEventsToResultEvent(events))

	// confirm that our memstore
	outputContract, err := proxy.MemStore.Get(inputContract.Key())
	require.NoError(t, err)
	require.Equal(t, inputContract, outputContract)

	// confirm that we can close the contract
	proxy.MemStore.SetHeight(200)
	inputContract.SettlementHeight = 200
	closeEvents := sdk.Events{
		types.NewCloseContractEvent(&inputContract),
	}
	proxy.handleCloseContractEvent(convertEventsToResultEvent(closeEvents))
	_, err = proxy.MemStore.Get(inputContract.Key()) // contract should be deleted from store since its closed
	require.Error(t, err)
}

func TestHandleHandleContractSettlementEvent(t *testing.T) {
	proxy := NewProxy(testConfig)
	inputContract := types.Contract{
		Provider:           testConfig.ProviderPubKey,
		Service:            common.BTCService,
		Client:             types.GetRandomPubKey(),
		Delegate:           common.EmptyPubKey,
		Type:               types.ContractType_PAY_AS_YOU_GO,
		Height:             100,
		Duration:           100,
		Rate:               cosmos.NewInt64Coin("uarkeo", 1),
		Deposit:            sdk.NewInt(100),
		Nonce:              0,
		Id:                 1,
		SettlementDuration: 10,
	}
	openCost := int64(100)
	events := sdk.Events{
		types.NewOpenContractEvent(openCost, &inputContract),
	}
	proxy.handleOpenContractEvent(convertEventsToResultEvent(events))

	// confirm that our memstore has the contract.
	outputContract, err := proxy.MemStore.Get(inputContract.Key())
	require.NoError(t, err)
	require.Equal(t, inputContract, outputContract)

	// simulate 10 calls being made to sentinel on the contract.
	arkAuth := ArkAuth{
		ContractId: inputContract.Id,
		Spender:    inputContract.Client,
		Nonce:      10,
	}
	_, err = proxy.paidTier(arkAuth, "")
	require.NoError(t, err)

	// confirm our claim exists in the claim store
	claim, err := proxy.ClaimStore.Get(Claim{ContractId: inputContract.Id}.Key())
	require.NoError(t, err)
	require.Equal(t, claim.ContractId, inputContract.Id)
	require.Equal(t, claim.Nonce, arkAuth.Nonce)

	// confirm is a settlement event is emitted with a lower nonce, we handle it correctly, by not setting our claim to Claimed.
	inputContract.Nonce = 8
	proxy.MemStore.SetHeight(150)
	settlementEvents := sdk.Events{
		types.NewContractSettlementEvent(sdk.NewInt(8), sdk.NewInt(1), &inputContract),
	}
	proxy.handleContractSettlementEvent(convertEventsToResultEvent(settlementEvents))

	claim, err = proxy.ClaimStore.Get(Claim{ContractId: inputContract.Id}.Key())
	require.NoError(t, err)
	require.Equal(t, claim.ContractId, inputContract.Id)
	require.Equal(t, claim.Nonce, arkAuth.Nonce)
	require.False(t, claim.Claimed)

	// confirm is a settlement event is emitted with the samce nonce, we handle it correctly, by setting our claim to Claimed.
	inputContract.Nonce = 10
	proxy.MemStore.SetHeight(160)
	settlementEvents = sdk.Events{
		types.NewContractSettlementEvent(sdk.NewInt(10), sdk.NewInt(1), &inputContract),
	}
	proxy.handleContractSettlementEvent(convertEventsToResultEvent(settlementEvents))
	claim, err = proxy.ClaimStore.Get(Claim{ContractId: inputContract.Id}.Key())
	require.NoError(t, err)
	require.Equal(t, claim.ContractId, inputContract.Id)
	require.Equal(t, claim.Nonce, arkAuth.Nonce)
	require.True(t, claim.Claimed)
}

func TestHandleNewBlockHeaderEvent(t *testing.T) {
	// TODO: add tests
}

func convertEventsToResultEvent(events sdk.Events) tmCoreTypes.ResultEvent {
	return tmCoreTypes.ResultEvent{
		Events: stringifyEvents(events.ToABCIEvents()),
	}
}

// stringifyEvents - adapated from the tendermint codebase
func stringifyEvents(events []abciTypes.Event) map[string][]string {
	result := make(map[string][]string)
	for _, event := range events {
		if len(event.Type) == 0 {
			continue
		}
		for _, attr := range event.Attributes {
			if len(attr.Key) == 0 {
				continue
			}
			compositeTag := fmt.Sprintf("%s.%s", event.Type, string(attr.Key))
			result[compositeTag] = append(result[compositeTag], string(attr.Value))
		}
	}
	return result
}
