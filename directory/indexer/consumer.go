package indexer

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/types"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

type attributes func() map[string]string

func parseEventToEventOpenContract(event interface{}) (atypes.EventOpenContract, error) {
	eventData := make(map[string]string)
	prefix := "arkeo.arkeo.EventOpenContract."
	switch evt := event.(type) {
	case ctypes.ResultEvent:
		for key, attribute := range evt.Events {
			k := strings.TrimPrefix(key, prefix)
			v := strings.Trim(attribute[0], `"`)
			eventData[k] = v
		}
	case abcitypes.Event:
		for _, attribute := range evt.Attributes {
			key := strings.TrimPrefix(string(attribute.GetKey()), prefix)
			value := strings.Trim(string(attribute.GetValue()), `"`)
			eventData[key] = value
		}
	default:
		return atypes.EventOpenContract{}, fmt.Errorf("unsupported event type: %T", evt)
	}

	type eventOpenContractAlias atypes.EventOpenContract
	eventOpenContract := struct {
		ContractId         string `json:"contract_id,omitempty"`
		Height             string `json:"height,omitempty"`
		Duration           string `json:"duration,omitempty"`
		Rate               string `json:"rate,omitempty"`
		Deposit            string `json:"deposit,omitempty"`
		Type               string `json:"type"`
		OpenCost           string `json:"open_cost"`
		SettlementDuration string `json:"settlement_duration"`
		Authorization      string `json:"authorization"`
		QueriesPerMinute   string `json:"queries_per_minute"`
		eventOpenContractAlias
	}{}

	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}

	if err := json.Unmarshal(jsonData, &eventOpenContract); err != nil {
		return atypes.EventOpenContract{}, err
	}

	result := atypes.EventOpenContract(eventOpenContract.eventOpenContractAlias)

	// make conversions
	result.QueriesPerMinute, err = strconv.ParseInt(eventOpenContract.QueriesPerMinute, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.OpenCost, err = strconv.ParseInt(eventOpenContract.OpenCost, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.Deposit, _ = cosmos.NewIntFromString(eventOpenContract.Deposit)
	result.ContractId, err = strconv.ParseUint(eventOpenContract.ContractId, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.Height, err = strconv.ParseInt(eventOpenContract.Height, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.SettlementDuration, err = strconv.ParseInt(eventOpenContract.SettlementDuration, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.Duration, err = strconv.ParseInt(eventOpenContract.Duration, 10, 64)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	err = json.Unmarshal([]byte(eventOpenContract.Rate), &result.Rate)
	if err != nil {
		return atypes.EventOpenContract{}, err
	}
	result.Authorization = atypes.ContractAuthorization(atypes.ContractAuthorization_value[eventOpenContract.Authorization])
	result.Type = atypes.ContractType(atypes.ContractType_value[eventOpenContract.Type])
	return result, nil
}

func parseEventToEventModProvider(event interface{}) (atypes.EventModProvider, error) {
	eventData := make(map[string]string)

	prefix := "arkeo.arkeo.EventModProvider."
	switch evt := event.(type) {
	case ctypes.ResultEvent:
		for key, attribute := range evt.Events {
			k := strings.TrimPrefix(key, prefix)
			v := strings.Trim(attribute[0], `"`)
			eventData[k] = v
		}
	case abcitypes.Event:
		for _, attribute := range evt.Attributes {
			key := strings.TrimPrefix(string(attribute.GetKey()), prefix)
			value := strings.Trim(string(attribute.GetValue()), `"`)
			eventData[key] = value
		}
	default:
		return atypes.EventModProvider{}, fmt.Errorf("unsupported event type: %T", evt)
	}

	type eventModProviderAlias atypes.EventModProvider
	eventModProvider := struct {
		MaxContractDuration string `json:"max_contract_duration,omitempty"`
		MinContractDuration string `json:"min_contract_duration,omitempty"`
		SettlementDuration  string `json:"settlement_duration,omitempty"`
		MetadataNonce       string `json:"metadata_nonce,omitempty"`
		SubscriptionRate    string `json:"subscription_rate"`
		PayAsYouGoRate      string `json:"pay_as_you_go_rate"`
		Status              string `json:"status"`
		eventModProviderAlias
	}{}

	jsonData, err := json.Marshal(eventData)
	if err != nil {
		return atypes.EventModProvider{}, err
	}

	if err := json.Unmarshal(jsonData, &eventModProvider); err != nil {
		return atypes.EventModProvider{}, err
	}

	result := atypes.EventModProvider(eventModProvider.eventModProviderAlias)

	// make conversions
	result.MaxContractDuration, err = strconv.ParseInt(eventModProvider.MaxContractDuration, 10, 64)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	result.MinContractDuration, err = strconv.ParseInt(eventModProvider.MinContractDuration, 10, 64)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	result.SettlementDuration, err = strconv.ParseInt(eventModProvider.SettlementDuration, 10, 64)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	result.MetadataNonce, err = strconv.ParseUint(eventModProvider.MetadataNonce, 10, 64)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	err = json.Unmarshal([]byte(eventModProvider.SubscriptionRate), &result.SubscriptionRate)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	err = json.Unmarshal([]byte(eventModProvider.PayAsYouGoRate), &result.PayAsYouGoRate)
	if err != nil {
		return atypes.EventModProvider{}, err
	}
	result.Status = atypes.ProviderStatus(atypes.ProviderStatus_value[eventModProvider.Status])
	return result, nil
}

func tmAttributeSource(tx tmtypes.Tx, evt abcitypes.Event, height int64) func() map[string]string {
	attribs := make(map[string]string, 0)
	for _, attr := range evt.Attributes {
		attribs[string(attr.Key)] = strings.Trim(string(attr.Value), `"`)
	}

	if tx != nil {
		if _, ok := attribs["hash"]; !ok {
			attribs["hash"] = strings.ToUpper(hex.EncodeToString(tx.Hash()))
		}
	}

	attribs["eventHeight"] = strconv.FormatInt(height, 10)
	if _, ok := attribs["height"]; !ok {
		attribs["height"] = attribs["eventHeight"]
	}

	return func() map[string]string { return attribs }
}

func (s *Service) handleValidatorPayoutEvent(evt types.ValidatorPayoutEvent) error {
	s.logger.Infof("received validatorPayoutEvent %#v", evt)
	if evt.Paid < 0 {
		return fmt.Errorf("received negative paid amt: %d for tx %s", evt.Paid, evt.TxID)
	}
	if evt.Paid == 0 {
		return nil
	}
	s.logger.Infof("upserting validator payout event for tx %s", evt.TxID)
	if _, err := s.db.UpsertValidatorPayoutEvent(evt); err != nil {
		return errors.Wrapf(err, "error upserting validator payout event")
	}
	return nil
}

// consumeEvents make connection to tendermint using websocket and then consume NewBlock event
func (s *Service) consumeEvents() error {
	s.logger.WithField("websocket", s.params.TendermintWs).Info("starting realtime indexing using /websocket")
	client, err := utils.NewTendermintClient(s.params.TendermintWs)
	if err != nil {
		return fmt.Errorf("fail to create tm client for %s, err: %w", s.params.TendermintWs, err)
	}

	if err = client.Start(); err != nil {
		return fmt.Errorf("fail to start websocket client,endpoint:%s,err: %w", s.params.TendermintWs, err)
	}
	defer func() {
		if err := client.Stop(); err != nil {
			s.logger.WithError(err).Error("error stopping client")
		}
	}()
	// splitting across multiple tendermint clients as websocket allows max of 5 subscriptions per client
	blockEvents, err := subscribe(context.Background(), client, "tm.event = 'NewBlock'")
	if err != nil {
		return err
	}

	s.logger.Info("beginning realtime event consumption")
	for {
		select {
		case evt := <-blockEvents:
			data, ok := evt.Data.(tmtypes.EventDataNewBlock)
			if !ok {
				continue
			}
			s.logger.WithField("height", data.Block.Height).Debug("received block")
			if err := s.gapFiller(); err != nil {
				s.logger.WithError(err).Error("fail to create block gap")
			}
		}
	}
}

// consumeHistoricalBlock index one block at a time
func (s *Service) consumeHistoricalBlock(blockHeight int64) (result *db.Block, err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var block *ctypes.ResultBlock
	var blockResults *ctypes.ResultBlockResults
	var blockErr, resultsErr error

	go func() {
		defer wg.Done()
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), defaultRetrieveBlockTimeout)
		defer cancel()
		block, blockErr = s.tmClient.Block(ctx, &blockHeight)
		// TODO: change this to use prometheus
		if time.Since(start) > 500*time.Millisecond {
			s.logger.Warnf("%.3f elapsed reading block %d", time.Since(start).Seconds(), blockHeight)
		}
	}()

	go func() {
		defer wg.Done()
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), defaultRetrieveBlockTimeout)
		defer cancel()
		blockResults, resultsErr = s.tmClient.BlockResults(ctx, &blockHeight)
		if time.Since(start) > 500*time.Millisecond {
			s.logger.Warnf("%.3f elapsed reading block results %d", time.Since(start).Seconds(), blockHeight)
		}
	}()
	wg.Wait()

	if blockErr != nil {
		return nil, fmt.Errorf("fail to read block,err: %w", blockErr)
	}
	if resultsErr != nil {
		return nil, fmt.Errorf("fail to read blockresult,err: %w", resultsErr)
	}

	log := s.logger.WithField("height", block.Block.Height)
	for _, transaction := range block.Block.Txs {
		if err := s.handleTransaction(block.Block.Height, transaction); err != nil {
			log.WithError(err).Error("fail to handler transaction")
		}
	}

	for _, event := range blockResults.EndBlockEvents {
		log.Debugf("received %s endblock event", event.Type)
		if err := s.handleAbciEvent(event, nil, block.Block.Height); err != nil {
			log.WithError(err).Errorf("error handling abci event %#v", event)
		}
	}

	r := &db.Block{
		Height:    block.Block.Height,
		Hash:      block.Block.Hash().String(),
		BlockTime: block.Block.Time,
	}
	return r, nil
}
func (s *Service) handleTransaction(height int64, transaction tmtypes.Tx) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultRetrieveTransactionTimeout)
	defer cancel()
	txInfo, err := s.tmClient.Tx(ctx, transaction.Hash(), false)
	if err != nil {
		return fmt.Errorf("failed to get transaction data for %s,err:%w", transaction.Hash(), err)
	}
	for _, event := range txInfo.TxResult.Events {
		s.logger.WithField("height", height).Debugf("received %s txevent", event.Type)
		if err := s.handleAbciEvent(event, transaction, height); err != nil {
			// move on
			s.logger.WithError(err).Errorf("error handling abci event %#v", event)
		}
	}
	return nil
}
func (s *Service) handleAbciEvent(event abcitypes.Event, transaction tmtypes.Tx, height int64) error {
	s.logger.WithField("height", height).
		WithField("type", event.Type).Info("handle abci event")
	switch event.Type {
	case atypes.EventTypeBondProvider:
		bondProviderEvent := types.BondProviderEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &bondProviderEvent); err != nil {
			return err
		}
		if err := s.handleBondProviderEvent(bondProviderEvent); err != nil {
			return err
		}
	case atypes.EventTypeModProvider:
		modProviderEvent, err := parseEventToEventModProvider(event)
		if err != nil {
			return err
		}
		if err := s.handleModProviderEvent(modProviderEvent); err != nil {
			return err
		}
	case atypes.EventTypeOpenContract:
		openContractEvent, err := parseEventToEventOpenContract(event)
		if err != nil {
			return err
		}
		if err := s.handleOpenContractEvent(openContractEvent); err != nil {
			return err
		}
	case atypes.EventTypeSettleContract:
		contractSettlementEvent := types.ContractSettlementEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &contractSettlementEvent); err != nil {
			return err
		}
		if err := s.handleContractSettlementEvent(contractSettlementEvent); err != nil {
			return err
		}
	case atypes.EventTypeValidatorPayout:
		validatorPayoutEvent := types.ValidatorPayoutEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &validatorPayoutEvent); err != nil {
			return err
		}
		if err := s.handleValidatorPayoutEvent(validatorPayoutEvent); err != nil {
			return err
		}
	case atypes.EventTypeCloseContract:
		s.logger.Debugf("received close_contract event")
		closeContractEvent := types.CloseContractEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &closeContractEvent); err != nil {
			return err
		}
		if err := s.handleCloseContractEvent(closeContractEvent); err != nil {
			return err
		}
	case "coin_spent", "coin_received", "transfer", "message", "tx":
		// do nothing
	default:
		// panic to make it immediately obvious that something is not handled
		// by directory indexer
		s.logger.Panicf("unrecognized event %s", event.Type)
	}
	return nil
}

// convertEvent copy attributes of map given by attributeFunc() to target which must be a pointer (map/slice implicitly ptr)
func convertEvent(attributeFunc attributes, target interface{}) error {
	return mapstructure.WeakDecode(attributeFunc(), target)
}

func subscribe(ctx context.Context, client *tmclient.HTTP, query string) (<-chan ctypes.ResultEvent, error) {
	out, err := client.Subscribe(ctx, "", query)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to query,query:%s,err: %w", query, err)
	}
	return out, nil
}
