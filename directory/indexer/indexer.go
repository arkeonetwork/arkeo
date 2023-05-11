package indexer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"

	// "github.com/arkeonetwork/common/logging"
	// arkutils "github.com/arkeonetwork/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
)

var log = logging.WithoutFields()

type IndexerAppParams struct {
	ArkeoApi            string
	TendermintApi       string
	TendermintWs        string
	ChainID             string
	Bech32PrefixAccAddr string
	Bech32PrefixAccPub  string
	IndexerID           int64
	db.DBConfig
}

type IndexerApp struct {
	Height         int64
	params         IndexerAppParams
	db             *db.DirectoryDB
	done           chan struct{}
	blockProcessor chan int64
	blockMutex     sync.Mutex
}

func NewIndexer(params IndexerAppParams) *IndexerApp {
	d, err := db.New(params.DBConfig)
	if err != nil {
		panic(fmt.Sprintf("error connecting to the db: %+v", err))
	}
	return &IndexerApp{
		params:         params,
		db:             d,
		blockProcessor: make(chan int64),
		blockMutex:     sync.Mutex{},
	}
}

func (a *IndexerApp) Run() (done <-chan struct{}, err error) {
	// initialize by reading all existing providers?
	a.done = make(chan struct{})
	a.realtime()
	return a.done, nil
}

func NewTenderm1intClient(baseURL string) (*tmclient.HTTP, error) {
	client, err := tmclient.New(baseURL, "/websocket")
	if err != nil {
		return nil, errors.Wrapf(err, "error creating websocket client")
	}
	logger := tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
	client.SetLogger(logger)

	return client, nil
}

func (a *IndexerApp) gapFiller() {
	a.blockMutex.Lock()
	defer a.blockMutex.Unlock()

	var err error
	tm, err := utils.NewTendermintClient(a.params.TendermintWs)
	if err != nil {
		log.Panicf("error creating gapFiller client: %+v", err)
	}

	latestStored, err := a.db.FindLatestBlock()
	if err != nil {
		log.Panicf("error finding latest stored block: %+v", err)
	}

	ctx := context.Background()
	latest, err := tm.Block(ctx, nil)
	if err != nil {
		log.Panicf("error finding latest block: %+v", err)
	}
	if latest.Block == nil {
		log.Errorf("latest block is nil, skipping")
		return
	}

	var todo db.BlockGap

	if latest.Block.Height-latestStored.Height <= 0 {
		return
	}
	if latestStored.Height == 0 {
		todo = db.BlockGap{Start: 1, End: latest.Block.Height}
	} else {
		log.Infof("%d missed blocks from %d to current %d", latest.Block.Height-latestStored.Height, latestStored.Height, latest.Block.Height)
		todo = db.BlockGap{Start: latestStored.Height + 1, End: latest.Block.Height}
	}

	if err := a.fillGap(todo); err != nil {
		log.Errorf("error filling gap %s", todo)
	}
}

// gaps filled inclusively
func (a *IndexerApp) fillGap(gap db.BlockGap) error {
	log.Infof("gap filling %s", gap)
	tm, err := utils.NewTendermintClient(a.params.TendermintWs)
	if err != nil {
		return errors.Wrapf(err, "error creating tm client: %+v", err)
	}

	for i := gap.Start; i <= gap.End; i++ {
		log.Infof("processing %d", i)
		// TODO: should pass in a db.Begin()/db.Commit() to ensure all or
		// nothing gets written
		block, err := a.consumeHistoricalBlock(tm, i)
		if err != nil {
			log.Errorf("error consuming block %d: %+v", i, err)
			continue
		}
		if _, err = a.db.InsertBlock(block); err != nil {
			log.Errorf("error inserting block %d with hash %s: %+v", block.Height, block.Hash, err)
			time.Sleep(time.Second)
		}
	}
	return nil
}

const numClients = 3

func (a *IndexerApp) realtime() {
	log.Infof("starting realtime indexing using /websocket at %s", a.params.TendermintWs)
	clients := make([]*tmclient.HTTP, numClients)
	for i := 0; i < numClients; i++ {
		client, err := utils.NewTendermintClient(a.params.TendermintWs)
		if err != nil {
			panic(fmt.Sprintf("error creating tm client for %s: %+v", a.params.TendermintWs, err))
		}
		if err = client.Start(); err != nil {
			panic(fmt.Sprintf("error starting ws client: %s: %+v", a.params.TendermintWs, err))
		}
		defer func() {
			if err := client.Stop(); err != nil {
				log.Errorf("error stopping client: %+v", err)
			}
		}()
		clients[i] = client
	}

	if err := a.consumeEvents(clients); err != nil {
		log.Errorf("error consuming events: %+v", err)
	}
	a.done <- struct{}{}
}
