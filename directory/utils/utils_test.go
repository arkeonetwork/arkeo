package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/arkeonetwork/arkeo/directory/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestParseCoordinates(t *testing.T) {
	epsilon := .0001
	coordinateString := "67.3523,-47.6878"
	coordinates, err := ParseCoordinates(coordinateString)
	if err != nil {
		t.FailNow()
	}
	if !IsNearEqual(coordinates.Latitude, 67.35234, epsilon) ||
		!IsNearEqual(coordinates.Longitude, -47.6878, epsilon) {
		t.FailNow()
	}

	coordinateString = "67.3523,-x"
	coordinates, err = ParseCoordinates(coordinateString)
	if err == nil {
		t.FailNow()
	}

	coordinateString = "yy,-47.6878"
	coordinates, err = ParseCoordinates(coordinateString)
	if err == nil {
		t.FailNow()
	}

	coordinateString = "67.3523,-47.6878,666"
	coordinates, err = ParseCoordinates(coordinateString)
	if err == nil {
		t.FailNow()
	}
	coordinateString = "67.3523"
	coordinates, err = ParseCoordinates(coordinateString)
	if err == nil {
		t.FailNow()
	}
}

func TestParseContractType(t *testing.T) {
	contract := "paygo"
	_, err := ParseContractType(contract)
	if err == nil {
		t.FailNow()
	}

	contract = "PayAsYouGo"
	contractType, err := ParseContractType(contract)
	if err != nil {
		t.FailNow()
	}

	if contractType != types.ContractTypePayAsYouGo {
		t.FailNow()
	}

	contract = "Subscription"
	contractType, err = ParseContractType(contract)
	if err != nil {
		t.FailNow()
	}
	if contractType != types.ContractTypeSubscription {
		t.FailNow()
	}
}

func TestDownloadProviderMetadata(t *testing.T) {
	sdkConfig := sdk.GetConfig()
	sdkConfig.SetBech32PrefixForAccount("tarkeo", "tarkeopub")

	// Open the sample file from the disk
	sourceFile, err := os.Open("../../testutil/sample/metadata.json")
	if err != nil {
		t.Fatalf("Failed to open source file: %v", err)
	}
	defer sourceFile.Close()

	// Create a temporary HTTP server to serve the file
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeContent(w, r, "metadata.json", time.Now(), sourceFile)
	}))
	defer ts.Close()

	metadata, err := DownloadProviderMetadata(ts.URL, 5, 1e6)
	if err != nil {
		t.FailNow()
	}

	if metadata == nil {
		t.FailNow()
	}

	if metadata.Version == "" {
		t.FailNow()
	}

	_, err = DownloadProviderMetadata(ts.URL, 5, 1)
	if err == nil {
		t.FailNow()
	}
}
