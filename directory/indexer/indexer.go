package indexer

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"

	// "github.com/arkeonetwork/common/logging"
	// arkutils "github.com/arkeonetwork/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	tmtypes "github.com/tendermint/tendermint/types"
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
	Height   int64
	IsSynced atomic.Bool
	params   IndexerAppParams
	db       *db.DirectoryDB
	done     chan struct{}
}

func NewIndexer(params IndexerAppParams) *IndexerApp {
	d, err := db.New(params.DBConfig)
	if err != nil {
		panic(fmt.Sprintf("error connecting to the db: %+v", err))
	}
	return &IndexerApp{params: params, db: d}
}

func (a *IndexerApp) Run() (done <-chan struct{}, err error) {
	// initialize by reading all existing providers?
	a.done = make(chan struct{})
	go a.realtime()
	go a.gapFiller()
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

const fillThreads = 3

var gaps []*db.BlockGap

func (a *IndexerApp) gapFiller() {
	var err error
	workChan := make(chan *db.BlockGap, 64)
	tm, err := utils.NewTendermintClient(a.params.TendermintWs)
	if err != nil {
		log.Panicf("error creating gapFiller client: %+v", err)
	}

	ctx := context.Background()
	for {
		gaps, err = a.db.FindBlockGaps()
		if err != nil {
			log.Errorf("error reading blocks from db: %+v", err)
		}

		latestStored, err := a.db.FindLatestBlock()
		if err != nil {
			log.Panicf("error finding latest stored block: %+v", err)
		}

		latest, err := tm.Block(ctx, nil)
		if err != nil {
			log.Panicf("error finding latest block: %+v", err)
		}

		if latestStored == nil {
			log.Infof("no latestStored, initializing")
			gaps = append(gaps, &db.BlockGap{Start: 1, End: latest.Block.Height})
		} else if latest.Block.Height-latestStored.Height > 1 {
			log.Infof("%d missed blocks from %d to current %d", latest.Block.Height-latestStored.Height, latestStored.Height, latest.Block.Height)
			gaps = append(gaps, &db.BlockGap{Start: latestStored.Height + 1, End: latest.Block.Height - 1})
		}

		if len(gaps) > 0 {
			log.Infof("have %d gaps to fill: %s", len(gaps), gaps)
			for i := range gaps {
				workChan <- gaps[i]
			}

			startThreads := len(gaps)
			if startThreads > fillThreads {
				startThreads = fillThreads
			}

			var wg sync.WaitGroup
			wg.Add(len(gaps))
			for i := 0; i < startThreads; i++ {
				go func() {
					for {
						select {
						case g := <-workChan:
							if err := a.fillGap(*g); err != nil {
								log.Errorf("error filling gap %s", g)
							}
							wg.Done()
						default:
							log.Infof("no work delivered, done")
							return
						}
					}
				}()
			}
			log.Infof("waiting for %d threads to complete filling %d gaps", startThreads, len(gaps))
			wg.Wait()
		}

		// all gaps filled, wait a minute
		time.Sleep(time.Minute)
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
		defer client.Stop()
		clients[i] = client
	}

	a.consumeEvents(clients)
	a.done <- struct{}{}
}

func (a *IndexerApp) handleBlockEvent(block *tmtypes.Block) error {
	if _, err := a.db.InsertBlock(&db.Block{Height: block.Height, Hash: block.Hash().String(), BlockTime: block.Time}); err != nil {
		return errors.Wrapf(err, "error inserting block")
	}
	a.Height = block.Height
	return nil
}
