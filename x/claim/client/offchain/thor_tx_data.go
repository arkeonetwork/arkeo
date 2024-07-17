package offchain

import (
	"fmt"
	"io"
	"net/http"

	"cosmossdk.io/errors"
)

func fetchThorChainTxData(hash string) ([]byte, error) {
	url := fmt.Sprintf("https://thornode.ninerealms.com/thorchain/tx/%s", hash)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build request %s", req.RequestURI)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get thorchain tx for %s", hash)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
