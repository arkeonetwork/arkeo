package indexer

import (
	"context"
	"fmt"
	"sync"
	"time"

	tmclient "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
)

const (
	defaultRetrieveBlockTimeout       = time.Second * 5
	defaultRetrieveTransactionTimeout = time.Second
)

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

// Service consume events from blockchain and persist it to a database
type Service struct {
	Height         int64
	params         ServiceParams
	db             *db.DirectoryDB
	done           chan struct{}
	blockProcessor chan int64
	blockMutex     sync.Mutex
	wg             *sync.WaitGroup
	logger         logging.Logger
	tmClient       *tmclient.HTTP
	blockFillQueue chan db.BlockGap
}

// NewIndexer create a new instance of Indexer
func NewIndexer(params ServiceParams) (*Service, error) {
	d, err := db.New(params.DB)
	if err != nil {
		return nil, fmt.Errorf("fail to connect to db,err: %w", err)
	}
	client, err := utils.NewTendermintClient(params.TendermintWs)
	if err != nil {
		return nil, fmt.Errorf("fail to create connection to tendermint,err:%w", err)
	}
	return &Service{
		params:         params,
		db:             d,
		blockProcessor: make(chan int64),
		blockMutex:     sync.Mutex{},
		done:           make(chan struct{}),
		logger: logging.WithFields(
			logging.Fields{
				"service": "indexer",
			}),
		tmClient:       client,
		wg:             &sync.WaitGroup{},
		blockFillQueue: make(chan db.BlockGap),
	}, nil
}

func (s *Service) Run() error {
	s.logger.Info("start to indexer service")
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.consumeEvents(); err != nil {
			s.logger.WithError(err).Error("fail to consume events")
		}
	}()
	s.wg.Add(1)
	go s.blockGapProcessor()
	return nil
}

func (s *Service) gapFiller() error {
	latestStored, err := s.db.FindLatestBlock()
	if err != nil {
		return fmt.Errorf("fail to find latest store block,err: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultRetrieveBlockTimeout)
	defer cancel()
	latest, err := s.tmClient.Block(ctx, nil)
	if err != nil {
		return fmt.Errorf("fail to find latest block,err: %w", err)
	}
	if latest.Block == nil {
		s.logger.Info("latest block is nil, skipping")
		return nil
	}
	var todo db.BlockGap
	if latest.Block.Height-latestStored.Height <= 0 {
		return nil
	}

	start := latestStored.Height + 1
	s.logger.Infof("%d missed blocks from %d to current %d", latest.Block.Height-latestStored.Height, latestStored.Height, latest.Block.Height)
	todo = db.BlockGap{Start: start, End: latest.Block.Height}
	select {
	// blockFillQueue is a blocking channel, only one gapfill task can be push into this channel when block gap processor is waiting on the other end
	// this ensures the service is only indexing one block at a time
	// the service should not index multiple blocks in parallel , it could cause data corrupt
	// especially when indexer start from scratch , while arkeo blockchain already have thousands blocks already
	// for example , contract opened on block 1 , and closed on block 10 , if the service process multiple blocks at the same time
	// it might process block 10 first , and block 1 later, which might consider contract still open
	case s.blockFillQueue <- todo:
	case <-s.done:
		return nil
	default:
		s.logger.Info("still processing previous block,skip")
	}

	return nil
}

func (s *Service) blockGapProcessor() {
	defer s.wg.Done()
	for {
		select {
		case blockGap, more := <-s.blockFillQueue:
			if !more {
				return
			}
			if err := s.fillGap(blockGap); err != nil {
				s.logger.WithError(err).
					WithField("start", blockGap.Start).
					WithField("end", blockGap.End).
					Error("fail to process blockgap")
			}
		case <-s.done:
			return
		}
	}
}

// gaps filled inclusively
func (s *Service) fillGap(gap db.BlockGap) error {
	s.logger.Infof("gap filling %s", gap)

	for i := gap.Start; i <= gap.End; i++ {
		s.logger.Infof("processing block: %d", i)
		block, err := s.consumeHistoricalBlock(i)
		if err != nil {
			s.logger.WithError(err).Errorf("err consuming block %d:", i)
			continue
		}
		if _, err = s.db.InsertBlock(block); err != nil {
			s.logger.WithError(err).Errorf("error inserting block %d with hash %s", block.Height, block.Hash)
		}
	}
	return nil
}

// Close will be called when it is time to shut down the service
// this allows the service to shut down itself gracefully
func (s *Service) Close() error {
	close(s.done)
	close(s.blockFillQueue)
	s.wg.Wait()
	return nil
}
