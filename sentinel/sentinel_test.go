package sentinel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestHandleActiveContract(t *testing.T) {
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

	url := fmt.Sprintf("%s%s/%s", RoutesActiveContract, inputContract.Service, inputContract.GetSpender())

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	response := httptest.NewRecorder()
	proxy.handleActiveContract(response, req)

	// Check the status code is what we expect.
	require.Equal(t, http.StatusOK, response.Code)

	// Check the response body is what we expect.
	expectedReturn, err := json.Marshal(inputContract)
	require.NoError(t, err)
	require.Equal(t, string(expectedReturn), response.Body.String())

	// confirm failure with incorrect path variables
	url = fmt.Sprintf("%s%s/%s", RoutesActiveContract, inputContract.GetSpender(), testConfig.ProviderPubKey)
	req, err = http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	response = httptest.NewRecorder()
	proxy.handleActiveContract(response, req)
	require.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandleClaim(t *testing.T) {
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
		Id:                 151,
		SettlementDuration: 10,
	}
	openCost := int64(100)
	openEvent := types.NewOpenContractEvent(openCost, &inputContract)
	sdkEvt, err := sdk.TypedEventToEvent(&openEvent)
	require.NoError(t, err)

	resultEvent := makeResultEvent(sdkEvt, openEvent.Height)
	proxy.handleOpenContractEvent(resultEvent)

	// simulate 10 calls being made to sentinel on the contract.
	arkAuth := ArkAuth{
		ContractId: inputContract.Id,
		Spender:    inputContract.Client,
		Nonce:      10,
	}
	_, err = proxy.paidTier(arkAuth, "")
	require.NoError(t, err)

	// get the expected claim
	claim := NewClaim(inputContract.Id, nil, 0, "")
	claim, err = proxy.ClaimStore.Get(claim.Key())
	require.NoError(t, err)
	require.False(t, claim.Claimed)
	require.Equal(t, inputContract.Id, claim.ContractId)

	// we should have a valid claim in our store.
	// check that we can retrieve it.
	url := fmt.Sprintf("%s%d", RoutesClaim, inputContract.Id)

	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	response := httptest.NewRecorder()
	proxy.handleClaim(response, req)

	// Check the status code is what we expect.
	require.Equal(t, http.StatusOK, response.Code)

	expectedReturn, err := json.Marshal(claim)
	require.NoError(t, err)
	require.Equal(t, string(expectedReturn), response.Body.String())

	// confirm failure with incorrect path variables
	url = RoutesClaim
	req, err = http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	response = httptest.NewRecorder()
	proxy.handleActiveContract(response, req)
	require.Equal(t, http.StatusBadRequest, response.Code)
}

func TestHandleOpenClaims(t *testing.T) {
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
		Id:                 151,
		SettlementDuration: 10,
	}
	openCost := int64(100)
	openEvent := types.NewOpenContractEvent(openCost, &inputContract)
	sdkEvt, err := sdk.TypedEventToEvent(&openEvent)
	require.NoError(t, err)

	resultEvent := makeResultEvent(sdkEvt, openEvent.Height)
	proxy.handleOpenContractEvent(resultEvent)

	// simulate 10 calls being made to sentinel on the contract.
	arkAuth := ArkAuth{
		ContractId: inputContract.Id,
		Spender:    inputContract.Client,
		Nonce:      10,
	}
	_, err = proxy.paidTier(arkAuth, "")
	require.NoError(t, err)

	// repeat for a second contract rom a different client
	inputContract.Client = types.GetRandomPubKey()
	inputContract.Id = 420
	openEvent = types.NewOpenContractEvent(openCost, &inputContract)
	sdkEvt, err = sdk.TypedEventToEvent(&openEvent)
	require.NoError(t, err)

	resultEvent = makeResultEvent(sdkEvt, openEvent.Height)
	proxy.handleOpenContractEvent(resultEvent)
	arkAuth = ArkAuth{
		ContractId: inputContract.Id,
		Spender:    inputContract.Client,
		Nonce:      15,
	}
	_, err = proxy.paidTier(arkAuth, "")
	require.NoError(t, err)

	// we should have 2 valid claim in our store.
	// check that we can retrieve them.
	url := RoutesOpenClaims
	req, err := http.NewRequest("GET", url, nil)
	require.NoError(t, err)
	response := httptest.NewRecorder()
	proxy.handleOpenClaims(response, req)

	// Check the status code is what we expect.
	require.Equal(t, http.StatusOK, response.Code)

	openClaims := make([]Claim, 0)
	require.NoError(t, json.Unmarshal(response.Body.Bytes(), &openClaims))
	require.Equal(t, 2, len(openClaims))
}
