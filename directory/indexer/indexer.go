package indexer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	tmclient "github.com/cometbft/cometbft/rpc/client/http"

	"github.com/arkeonetwork/arkeo/common/logging"
	"github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
)

const (
	defaultRetrieveBlockTimeout       = time.Second * 5
	defaultRetrieveTransactionTimeout = time.Second
	defaultHandleEventTimeout         = time.Second * 5
	defaultFindLastBlockTimeout       = time.Second
)

// Service consume events from blockchain and persist it to a database
type Service struct {
	params         ServiceParams
	db             db.IDataStorage
	done           chan struct{}
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
		params: params,
		db:     d,
		done:   make(chan struct{}),
		logger: logging.WithFields(
			logging.Fields{
				"service": "indexer",
			}),
		tmClient:       client,
		wg:             &sync.WaitGroup{},
		blockFillQueue: make(chan db.BlockGap),
	}, nil
}

// Run start the indexer service
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
	ctxFindLastBlock, cancelFindLastBlock := context.WithTimeout(context.Background(), defaultFindLastBlockTimeout)
	defer cancelFindLastBlock()
	latestStored, err := s.db.FindLatestBlock(ctxFindLastBlock)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return fmt.Errorf("fail to find latest store block,err: %w", err)
		}
		// start from zero
		latestStored = &db.Block{}
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultRetrieveBlockTimeout)
	defer cancel()
	latest, err := s.tmClient.Block(ctx, nil)
	if err != nil {
		return fmt.Errorf("fail to find latest block,err: %w", err)
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

// blockGapProcessor will be run in a separate go routine, it populates blockgap task from blockFillQueue
// and then process it block by block , it can only process one blockgap task at a time
// it will only exit when the service receive a SIGTERM
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

// fillGap will consume all the blocks from arkeo node, and index all the events in it
func (s *Service) fillGap(gap db.BlockGap) error {
	s.logger.Infof("gap filling %s", gap)

	for i := gap.Start; i <= gap.End; i++ {
		s.logger.Infof("processing block: %d", i)

		// Upsert indexer status to track latest block height
		if _, err := s.db.UpsertIndexerStatus(context.Background(), i); err != nil {
			s.logger.WithError(err).Errorf("failed to upsert indexer status for height %d", i)
		}

		block, err := s.consumeHistoricalBlock(i)
		if err != nil {
			s.logger.WithError(err).Errorf("err consuming block %d:", i)
			continue
		}
		if _, err = s.db.InsertBlock(context.Background(), block); err != nil {
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

// Stringfy marshal the given input into json
func Stringfy(input any) string {
	buf, err := json.Marshal(input)
	if err != nil {
		return fmt.Sprintf("fail to stringfy object,err: %s", err)
	}
	return string(buf)
}
