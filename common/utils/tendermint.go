package utils

import (
	"os"

	tmlog "github.com/cometbft/cometbft/libs/log"
	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	"github.com/pkg/errors"
)

func NewTendermintClient(baseURL string) (*tmclient.HTTP, error) {
	client, err := tmclient.New(baseURL, "/websocket")
	if err != nil {
		return nil, errors.Wrapf(err, "error creating websocket client")
	}
	logger := tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
	client.SetLogger(logger)

	return client, nil
}
