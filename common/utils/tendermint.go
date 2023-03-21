package utils

import (
	"os"

	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
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
