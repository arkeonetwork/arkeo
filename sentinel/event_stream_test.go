package sentinel

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/arkeonetwork/arkeo/app"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	abciTypes "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/proto/tendermint/version"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func newTestConfig() conf.Configuration {
	c := cosmos.GetConfig()
	c.SetBech32PrefixForAccount(app.AccountAddressPrefix, app.AccountAddressPrefix+"pub")
	return conf.Configuration{
		Moniker:            "Testy McTestface",
		Website:            "testing.com",
		Description:        "the best testnet ever",
		Location:           "100,100",
		Port:               "3636",
		SourceChain:        "http://localhost:1317", // this should point to arkeo rpc endpoints, but we can ignore for testing
		EventStreamHost:    "localhost",
		ProviderPubKey:     types.GetRandomPubKey(),
		FreeTierRateLimit:  100,
		ClaimStoreLocation: "",
	}
}

func TestHandleOpenContractEvent(t *testing.T) {
	testConfig := newTestConfig()
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
		QueriesPerMinute:   1,
	}
	openCost := int64(100)

	openEvent := types.NewOpenContractEvent(openCost, &inputContract)
	sdkEvt, err := sdk.TypedEventToEvent(&openEvent)
	require.NoError(t, err)

	resultEvent := makeResultEvent(sdkEvt, openEvent.Height)
	proxy.handleOpenContractEvent(resultEvent)

	// confirm that our memstore has the contract and its active
	outputContract, err := proxy.MemStore.Get(inputContract.Key())
	require.NoError(t, err)
	require.Equal(t, inputContract, outputContract)
	_, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Service, inputContract.Client)
	require.NoError(t, err)

	// confirm that a contract for a different provider doesn't get stored.
	inputContract.Provider = types.GetRandomPubKey()
	inputContract.Id = 2
	openEvent.Provider = inputContract.Provider
	openEvent.ContractId = inputContract.Id

	sdkEvt, err = sdk.TypedEventToEvent(&openEvent)
	require.NoError(t, err)

	resultEvent = makeResultEvent(sdkEvt, openEvent.Height)
	proxy.handleOpenContractEvent(resultEvent)
	_, err = proxy.MemStore.Get(inputContract.Key())
	require.Error(t, err)

	// confirm that we return the correct active contract when multiple contracts are present.
	proxy.MemStore.SetHeight(201) // contract with id 1 should now be expired

	_, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Service, inputContract.Client)
	require.Error(t, err) // contract not found since its expired.

	inputContract.Provider = testConfig.ProviderPubKey
	inputContract.Id = 3
	inputContract.Height = 200
	openEvent.Provider = inputContract.Provider
	openEvent.ContractId = inputContract.Id
	openEvent.Height = inputContract.Height

	sdkEvt, err = sdk.TypedEventToEvent(&openEvent)
	require.NoError(t, err)

	resultEvent = makeResultEvent(sdkEvt, openEvent.Height)
	proxy.handleOpenContractEvent(resultEvent)

	outputContract, err = proxy.MemStore.GetActiveContract(inputContract.Provider, inputContract.Service, inputContract.Client)
	require.NoError(t, err)
	require.Equal(t, inputContract, outputContract)
}

func TestHandleCloseContractEvent(t *testing.T) {
	testConfig := newTestConfig()
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
	openEvent := types.NewOpenContractEvent(openCost, &inputContract)
	sdkEvt, err := sdk.TypedEventToEvent(&openEvent)
	require.NoError(t, err)

	resultEvent := makeResultEvent(sdkEvt, openEvent.Height)
	proxy.handleOpenContractEvent(resultEvent)

	// confirm that our memstore
	outputContract, err := proxy.MemStore.Get(inputContract.Key())
	require.NoError(t, err)
	require.Equal(t, inputContract, outputContract)

	// confirm that we can close the contract
	proxy.MemStore.SetHeight(200)
	inputContract.SettlementHeight = 200
	closeEvent := types.EventCloseContract{
		ContractId: inputContract.Id,
		Provider:   inputContract.Provider,
		Service:    inputContract.Service.String(),
		Client:     inputContract.Client,
		Delegate:   inputContract.Delegate,
	}

	sdkEvt, err = sdk.TypedEventToEvent(&closeEvent)
	require.NoError(t, err)

	resultEvent = makeResultEvent(sdkEvt, inputContract.SettlementHeight)
	proxy.handleCloseContractEvent(resultEvent)
	_, err = proxy.MemStore.Get(inputContract.Key()) // contract should be deleted from store since its closed
	require.Error(t, err)
}

func TestHandleHandleContractSettlementEvent(t *testing.T) {
	testConfig := newTestConfig()
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
		QueriesPerMinute:   1,
	}
	openCost := int64(100)
	openEvent := types.NewOpenContractEvent(openCost, &inputContract)
	sdkEvt, err := sdk.TypedEventToEvent(&openEvent)
	require.NoError(t, err)
	resultEvent := makeResultEvent(sdkEvt, openEvent.Height)
	proxy.handleOpenContractEvent(resultEvent)

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
	settlementEvent := types.NewContractSettlementEvent(sdk.NewInt(8), sdk.NewInt(1), &inputContract)
	sdkEvt, err = sdk.TypedEventToEvent(&settlementEvent)
	require.NoError(t, err)

	resultEvent = makeResultEvent(sdkEvt, 151)
	proxy.handleContractSettlementEvent(resultEvent)

	claim, err = proxy.ClaimStore.Get(Claim{ContractId: inputContract.Id}.Key())
	require.NoError(t, err)
	require.Equal(t, claim.ContractId, inputContract.Id)
	require.Equal(t, claim.Nonce, arkAuth.Nonce)
	require.False(t, claim.Claimed)

	// confirm is a settlement event is emitted with the samce nonce, we handle it correctly, by setting our claim to Claimed.
	inputContract.Nonce = 10
	proxy.MemStore.SetHeight(160)
	settlementEvent = types.NewContractSettlementEvent(sdk.NewInt(10), sdk.NewInt(1), &inputContract)
	sdkEvt, err = sdk.TypedEventToEvent(&settlementEvent)
	require.NoError(t, err)

	resultEvent = makeResultEvent(sdkEvt, 161)
	proxy.handleContractSettlementEvent(resultEvent)
	claim, err = proxy.ClaimStore.Get(Claim{ContractId: inputContract.Id}.Key())
	require.NoError(t, err)
	require.Equal(t, claim.ContractId, inputContract.Id)
	require.Equal(t, claim.Nonce, arkAuth.Nonce)
	require.True(t, claim.Claimed)
}

func TestHandleNewBlockHeaderEvent(t *testing.T) {
	// TODO: add tests
}

func makeResultEvent(sdkEvent sdk.Event, height int64) tmCoreTypes.ResultEvent {
	evts := make(map[string][]string, len(sdkEvent.Attributes))
	for _, attr := range sdkEvent.Attributes {
		evts[string(attr.Key)] = []string{string(attr.Value)}
	}

	abciEvents := []abciTypes.Event{{
		Type:       sdkEvent.Type,
		Attributes: sdkEvent.Attributes,
	}}

	query := fmt.Sprintf("tm.event = 'Tx' AND message.action='/%s'", sdkEvent.Type)
	return tmCoreTypes.ResultEvent{
		Query:  query,
		Events: evts,
		Data: tmtypes.EventDataTx{
			TxResult: abciTypes.TxResult{
				Height: height,
				Index:  0,
				Tx:     []byte{},
				Result: abciTypes.ResponseDeliverTx{
					Events: abciEvents,
				},
			},
		},
	}
}

func createMsgOpenContractEvent() sdk.Event {

	contract := createTestContract()
	openCost := int64(100)

	event := types.NewOpenContractEvent(openCost, &contract)

	sdkEvent, err := sdk.TypedEventToEvent(&event)
	if err != nil {
		panic(err)
	}

	return sdkEvent
}

func createTestContract() types.Contract {
	return types.Contract{
		Provider:           types.GetRandomPubKey(),
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
		QueriesPerMinute:   1,
	}
}

func createMsgCloseContractEvent() sdk.Event {
	contract := createTestContract()
	event := types.NewCloseContractEvent(&contract)

	sdkEvent, err := sdk.TypedEventToEvent(&event)
	if err != nil {
		panic(err)
	}

	return sdkEvent

}

func createMsgClaimContractEvent() sdk.Event {
	contract := createTestContract()
	debts := sdk.NewInt(10)
	valIncome := sdk.NewInt(5)
	event := types.NewContractSettlementEvent(debts, valIncome, &contract)

	sdkEvent, err := sdk.TypedEventToEvent(&event)
	if err != nil {
		panic(err)
	}

	return sdkEvent
}

func createBlockHeaderEvent() sdk.Event {
	header := &tmtypes.Header{}
	header.Populate(
		version.Consensus{
			Block: 1,
			App:   1,
		},
		"test-chain",
		time.Now(),
		tmtypes.BlockID{
			Hash: []byte("block_hash"),
			PartSetHeader: tmtypes.PartSetHeader{
				Total: 1,
				Hash:  []byte("parts_byte"),
			},
		},
		[]byte("validator_hash"),
		[]byte("next_validator_hash"),
		[]byte("consensus_hash"),
		[]byte("app_hash"),
		[]byte("last_result_hash"),
		[]byte("proposer_address"),
	)
	event := tmtypes.EventDataNewBlockHeader{
		Header:           *header,
		NumTxs:           1,
		ResultBeginBlock: abciTypes.ResponseBeginBlock{},
		ResultEndBlock:   abciTypes.ResponseEndBlock{},
	}

	eventBytes, _ := json.Marshal(event)

	return sdk.Event{
		Type: tmtypes.EventNewBlockHeader,
		Attributes: []abciTypes.EventAttribute{
			{
				Key:   []byte("type"),
				Value: []byte("NewBlockHeadr"),
			},
			{
				Key:   []byte("data"),
				Value: eventBytes,
			},
		},
	}
}

func TestEventListenerOrderProcessing(t *testing.T) {

	testconfig := newTestConfig()
	proxy := NewProxy(testconfig)

	// blockEvent := makeResultEvent(createBlockHeaderEvent(), 1)
	// openContractEvent := makeResultEvent(createMsgOpenContractEvent(), 2)
	// closeContractEvent := makeResultEvent(createMsgCloseContractEvent(), 3)
	// settleContractClaim := makeResultEvent(createMsgClaimContractEvent(), 4)

	// receivedEvents := make(chan tmCoreTypes.ResultEvent, 4)
	// mockEventHandler := func(events tmCoreTypes.ResultEvent) {
	// 	receivedEvents <- events
	// }

	go proxy.Run()

}
