package offchain

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cosmossdk.io/errors"
)

type Coin struct {
	Asset  string `json:"asset"`
	Amount string `json:"amount"`
}

type Tx struct {
	ID          string      `json:"id"`
	Chain       string      `json:"chain"`
	FromAddress string      `json:"from_address"`
	ToAddress   string      `json:"to_address"`
	Coins       []Coin      `json:"coins"`
	Gas         interface{} `json:"gas"`
	Memo        string      `json:"memo"`
}

type ObservedTx struct {
	Tx Tx `json:"tx"`
}

type KeysignMetric struct {
	TxID         string      `json:"tx_id"`
	NodeTSSTimes interface{} `json:"node_tss_times"`
}

type ThorChainTxData struct {
	ObservedTx      ObservedTx    `json:"observed_tx"`
	ConsensusHeight int           `json:"consensus_height"`
	FinalisedHeight int           `json:"finalised_height"`
	KeysignMetric   KeysignMetric `json:"keysign_metric"`
}

type ThorTxData struct {
	Hash   string `json:"hash"`
	TxData string `json:"tx_data"`
}

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

	var result ThorChainTxData
	if err = json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("error unmarshalling response: %w", err)
	}

	resultBytes, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("error marshalling result: %w", err)
	}

	txDataHash := sha512.Sum512(resultBytes)
	txHex := base64.RawURLEncoding.EncodeToString(txDataHash[:])
	txDataHex := base64.RawURLEncoding.EncodeToString(resultBytes)

	txData := ThorTxData{
		Hash:   txHex,
		TxData: txDataHex,
	}

	txDataBytes, err := json.Marshal(txData)
	if err != nil {
		return "", fmt.Errorf("error marshalling txData: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(txDataBytes), nil
}
