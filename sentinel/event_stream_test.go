package sentinel

import (
	"fmt"
	"testing"
	"time"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
)

func TestHandleOpenContractEvent(t *testing.T) {
	config := conf.Configuration{
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
	proxy := NewProxy(config)
	inputContract := types.Contract{
		Provider:           config.ProviderPubKey,
		Chain:              common.BTCChain,
		Client:             types.GetRandomPubKey(),
		Delegate:           common.EmptyPubKey,
		Type:               types.ContractType_PAY_AS_YOU_GO,
		Height:             100,
		Duration:           100,
		Rate:               1,
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
	_, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Chain, inputContract.Client)
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

	_, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Chain, inputContract.Client)
	require.Error(t, err) // contract not found since its expired.

	inputContract.Provider = config.ProviderPubKey
	inputContract.Id = 3
	inputContract.Height = 200
	events = sdk.Events{
		types.NewOpenContractEvent(openCost, &inputContract),
	}
	proxy.handleOpenContractEvent(convertEventsToResultEvent(events))

	outputContract, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Chain, inputContract.Client)
	require.NoError(t, err)
	require.Equal(t, inputContract, outputContract)
}

func TestHandleCloseContractEvent(t *testing.T) {

}

func TestHandleClaimContractIncomeEevent(t *testing.T) {

}

func TestHandleNewBlockHeaderEvent(t *testing.T) {

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
