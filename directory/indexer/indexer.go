package indexer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"

	// "github.com/arkeonetwork/common/logging"
	// arkutils "github.com/arkeonetwork/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/pkg/errors"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
)

var log = logging.WithoutFields()

// ServiceParams hold all necessary parameters for indexer app to run
type ServiceParams struct {
	ArkeoApi            string      `mapstructure:"arkeo_api" json:"arkeo_api"`
	TendermintApi       string      `mapstructure:"tendermint_api" json:"tendermint_api"`
	TendermintWs        string      `mapstructure:"tendermint_ws" json:"tendermint_ws"`
	ChainID             string      `mapstructure:"chain_id" json:"chain_id"`
	Bech32PrefixAccAddr string      `mapstructure:"bech32_pref_acc_addr" json:"bech32_pref_acc_addr"`
	Bech32PrefixAccPub  string      `mapstructure:"bech32_pref_acc_pub" json:"bech32_pref_acc_pub"`
	IndexerID           int64       `json:"-"`
	DB                  db.DBConfig `mapstructure:"db" json:"db"`
}

// Service consume events from blockchain and persist it to  database
type Service struct {
	Height         int64
	params         ServiceParams
	db             *db.DirectoryDB
	done           chan struct{}
	blockProcessor chan int64
	blockMutex     sync.Mutex
}

func NewIndexer(params ServiceParams) *Service {
	d, err := db.New(params.DB)
	if err != nil {
		panic(fmt.Sprintf("error connecting to the db: %+v", err))
	}
	return &Service{
		params:         params,
		db:             d,
		blockProcessor: make(chan int64),
		blockMutex:     sync.Mutex{},
	}
}

func (s *Service) Run() (done <-chan struct{}, err error) {
	// initialize by reading all existing providers?
	s.done = make(chan struct{})
	s.realtime()
	return s.done, nil
}

func (s *Service) gapFiller() {
	s.blockMutex.Lock()
	defer s.blockMutex.Unlock()

	var err error
	tm, err := utils.NewTendermintClient(s.params.TendermintWs)
	if err != nil {
		log.Panicf("error creating gapFiller client: %+v", err)
	}

	latestStored, err := s.db.FindLatestBlock()
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

	if err := s.fillGap(todo); err != nil {
		log.Errorf("error filling gap %s", todo)
	}
}

// gaps filled inclusively
func (s *Service) fillGap(gap db.BlockGap) error {
	log.Infof("gap filling %s", gap)
	tm, err := utils.NewTendermintClient(s.params.TendermintWs)
	if err != nil {
		return errors.Wrapf(err, "error creating tm client: %+v", err)
	}

	for i := gap.Start; i <= gap.End; i++ {
		log.Infof("processing %d", i)
		// TODO: should pass in s db.Begin()/db.Commit() to ensure all or
		// nothing gets written
		block, err := s.consumeHistoricalBlock(tm, i)
		if err != nil {
			log.Errorf("error consuming block %d: %+v", i, err)
			continue
		}
		if _, err = s.db.InsertBlock(block); err != nil {
			log.Errorf("error inserting block %d with hash %s: %+v", block.Height, block.Hash, err)
			time.Sleep(time.Second)
		}
	}
	return nil
}

const numClients = 3

func (s *Service) realtime() {
	log.Infof("starting realtime indexing using /websocket at %s", s.params.TendermintWs)
	clients := make([]*tmclient.HTTP, numClients)
	for i := 0; i < numClients; i++ {
		client, err := utils.NewTendermintClient(s.params.TendermintWs)
		if err != nil {
			panic(fmt.Sprintf("error creating tm client for %s: %+v", s.params.TendermintWs, err))
		}
		if err = client.Start(); err != nil {
			panic(fmt.Sprintf("error starting ws client: %s: %+v", s.params.TendermintWs, err))
		}
		defer func() {
			if err := client.Stop(); err != nil {
				log.Errorf("error stopping client: %+v", err)
			}
		}()
		clients[i] = client
	}

	if err := s.consumeEvents(clients); err != nil {
		log.Errorf("error consuming events: %+v", err)
	}
	s.done <- struct{}{}
}
