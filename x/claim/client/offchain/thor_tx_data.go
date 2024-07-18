package offchain

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cosmossdk.io/errors"
	"github.com/arkeonetwork/arkeo/x/claim/types"
)

func fetchThorChainTxData(hash string) (string, error) {
	url := fmt.Sprintf("https://thornode.ninerealms.com/thorchain/tx/%s", hash)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.Wrap(err, "failed to build request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get thorchain tx for %s", hash)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	var result types.ThorChainTxData
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("error marshalling result: %w", err)
	}

	txDataHash := sha512.Sum512(resultBytes)
	txHashBase64 := base64.StdEncoding.EncodeToString(txDataHash[:])
	txDataBase64 := base64.StdEncoding.EncodeToString(resultBytes)

	txData := types.ThorTxData{
		Hash:   txHashBase64,
		TxData: txDataBase64,
	}

	txDataBytes, err := json.Marshal(txData)
	if err != nil {
		return "", fmt.Errorf("error marshalling txData: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(txDataBytes), nil
}
