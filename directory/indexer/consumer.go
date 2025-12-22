package indexer

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	abcitypes "github.com/cometbft/cometbft/abci/types"
	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	ctypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"

	"github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

const (
	defaultGapFillInterval    = 30 * time.Second
	defaultTMReconnectBackoff = 5 * time.Second
)

var warnedUnhandledAbciEventTypes sync.Map

// consumeEvents make connection to tendermint using websocket and then consume NewBlock event
func (s *Service) consumeEvents() error {
	s.logger.WithField("websocket", s.params.TendermintWs).Info("starting realtime indexing using /websocket")

	ticker := time.NewTicker(defaultGapFillInterval)
	defer ticker.Stop()

	for {
		if err := s.gapFiller(); err != nil {
			s.logger.WithError(err).Error("fail to create block gap")
		}

		client, err := utils.NewTendermintClient(s.params.TendermintWs)
		if err != nil {
			s.logger.WithError(err).WithField("endpoint", s.params.TendermintWs).Error("failed to create tendermint client")
			select {
			case <-time.After(defaultTMReconnectBackoff):
				continue
			case <-s.done:
				return nil
			}
		}

		if err := client.Start(); err != nil {
			s.logger.WithError(err).WithField("endpoint", s.params.TendermintWs).Error("failed to start websocket client")
			_ = client.Stop()
			select {
			case <-time.After(defaultTMReconnectBackoff):
				continue
			case <-s.done:
				return nil
			}
		}

		subCtx, cancel := context.WithCancel(context.Background())
		// splitting across multiple tendermint clients as websocket allows max of 5 subscriptions per client
		blockEvents, err := subscribe(subCtx, client, "tm.event = 'NewBlock'")
		if err != nil {
			s.logger.WithError(err).Error("failed to subscribe to NewBlock")
			cancel()
			_ = client.Stop()
			select {
			case <-time.After(defaultTMReconnectBackoff):
				continue
			case <-s.done:
				return nil
			}
		}

		s.logger.Info("beginning realtime event consumption")

	readLoop:
		for {
			select {
			case evt, ok := <-blockEvents:
				if !ok {
					s.logger.Warn("tendermint subscription closed; reconnecting")
					break readLoop
				}
				data, ok := evt.Data.(tmtypes.EventDataNewBlock)
				if !ok {
					continue
				}
				s.logger.WithField("height", data.Block.Height).Debug("received block")
				if err := s.gapFiller(); err != nil {
					s.logger.WithError(err).Error("fail to create block gap")
				}
			case <-ticker.C:
				if err := s.gapFiller(); err != nil {
					s.logger.WithError(err).Error("fail to create block gap")
				}
			case <-s.done: // finished
				cancel()
				if err := client.Stop(); err != nil {
					s.logger.WithError(err).Error("error stopping client")
				}
				return nil
			}
		}

		cancel()
		if err := client.Stop(); err != nil {
			s.logger.WithError(err).Error("error stopping client")
		}

		select {
		case <-time.After(defaultTMReconnectBackoff):
			continue
		case <-s.done:
			return nil
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

	if len(blockResults.TxsResults) != len(block.Block.Txs) {
		log.WithField("txs", len(block.Block.Txs)).
			WithField("tx_results", len(blockResults.TxsResults)).
			Warn("tx results count mismatch")
	}

	for i, transaction := range block.Block.Txs {
		var txResult *abcitypes.ExecTxResult
		if i < len(blockResults.TxsResults) {
			txResult = blockResults.TxsResults[i]
		}
		if err := s.handleTransaction(block.Block.Height, transaction, txResult); err != nil {
			log.WithError(err).Error("fail to handle transaction")
		}
	}

	for _, event := range blockResults.FinalizeBlockEvents {
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

func (s *Service) handleTransaction(height int64, transaction tmtypes.Tx, txResult *abcitypes.ExecTxResult) error {
	if txResult == nil {
		return fmt.Errorf("nil tx result for tx %s", hex.EncodeToString(transaction.Hash()))
	}
	for _, event := range txResult.Events {
		s.logger.WithField("height", height).Debugf("received %s txevent", event.Type)
		if err := s.handleAbciEvent(event, transaction, height); err != nil {
			// move on
			s.logger.WithError(err).WithField("event", Stringfy(event)).Errorf("error handling abci event")
		}
	}
	return nil
}

func (s *Service) handleAbciEvent(event abcitypes.Event, transaction tmtypes.Tx, height int64) error {
	var txID string
	if transaction != nil {
		txID = hex.EncodeToString(transaction.Hash())
	}
	s.logger.WithField("height", height).
		WithField("type", event.Type).
		WithField("tx_id", txID).
		Debug("handle abci event")

	ctx, cancel := context.WithTimeout(context.Background(), defaultHandleEventTimeout)
	defer cancel()
	switch event.Type {
	case atypes.EventTypeBondProvider:
		bondProviderEvent, err := parseEventToConcreteType[atypes.EventBondProvider](event)
		if err != nil {
			return err
		}
		if err := s.handleBondProviderEvent(ctx, bondProviderEvent, txID, height); err != nil {
			return err
		}
	case atypes.EventTypeModProvider:
		modProviderEvent, err := parseEventToConcreteType[atypes.EventModProvider](event)
		if err != nil {
			return err
		}
		if err := s.handleModProviderEvent(ctx, modProviderEvent, txID, height); err != nil {
			return err
		}
	case atypes.EventTypeOpenContract:
		contractOpenEvent, err := parseEventToConcreteType[atypes.EventOpenContract](event)
		if err != nil {
			return err
		}
		if err := s.handleOpenContractEvent(ctx, contractOpenEvent, txID, height); err != nil {
			return err
		}
	case atypes.EventTypeSettleContract:
		eventSettleContract, err := parseEventToConcreteType[atypes.EventSettleContract](event)
		if err != nil {
			return err
		}
		if err := s.handleContractSettlementEvent(ctx, eventSettleContract, txID, height); err != nil {
			return err
		}
	case atypes.EventTypeValidatorPayout:
		// Intentionally ignored: writing these to `validator_payout_events` is very high-volume
		// and isn't currently used by the directory/indexer.
	case atypes.EventTypeCloseContract:
		eventCloseContract, err := parseEventToConcreteType[atypes.EventCloseContract](event)
		if err != nil {
			return err
		}
		if err := s.handleCloseContractEvent(ctx, eventCloseContract, txID, height); err != nil {
			return err
		}
	// Proposal events
	case "submit_proposal", "proposal_deposit", "proposal_vote", "active_proposal", "inactive_proposal", "proposal_execution_failed":
		attrJSON, err := json.Marshal(event.Attributes)
		if err != nil {
			return err
		}
		if err := s.handleGenericEvent(ctx, event.Type, txID, height, attrJSON); err != nil {
			return err
		}
	// Staking module events
	case "create_validator", "edit_validator", "redelegate", "unbond", "complete_redelegation", "complete_unbonding", "begin_redelegate":
		attrJSON, err := json.Marshal(event.Attributes)
		if err != nil {
			return err
		}
		if err := s.handleGenericEvent(ctx, event.Type, txID, height, attrJSON); err != nil {
			return err
		}
	// Burn events
	case "burn", "slash", "jail", "unjail", "cosmos.authz.v1beta1.EventGrant", "cosmos.authz.v1beta1.EventRevoke", "set_withdraw_address", "withdraw_validator_commission":
		attrJSON, err := json.Marshal(event.Attributes)
		if err != nil {
			return err
		}
		if err := s.handleGenericEvent(ctx, event.Type, txID, height, attrJSON); err != nil {
			return err
		}
	// IBC Events
	case "channel_open_init", "channel_open_ack", "connection_open_init", "connection_open_ack", "create_client":
		attrJSON, err := json.Marshal(event.Attributes)
		if err != nil {
			return err
		}
		if err := s.handleGenericEvent(ctx, event.Type, txID, height, attrJSON); err != nil {
			return err
		}
	// IBC Events (Untracked)
	case "send_packet", "recv_packet", "acknowledge_packet", "write_acknowledgement", "timeout_packet", "channel_open_try", "channel_open_confirm", "channel_close_init", "channel_close_confirm", "connection_open_try", "connection_open_confirm", "ics20_transfer", "update_client", "timeout", "fungible_token_packet", "denomination_trace", "send_native", "recv_native", "timeout_native":
		// Not logged, too high volume.
	// Bridge/claim events
	case "withdraw_rewards", "delegate", "claim", "withdraw_commission", "claim_from_eth", "claim_thor_delegate":
		// Not logging these because they are tied to claims and rewards, but are a lot of data.
		// High volume, generally only useful for low-level validator monitoring.
	// Authz/Group Module
	case "cosmos.group.v1.EventCreateGroup", "cosmos.group.v1.EventExec", "cosmos.group.v1.EventLeaveGroup":
	// Extra Modules
	case "cosmos.feegrant.v1beta1.EventGrant", "cosmos.feegrant.v1beta1.EventRevoke", "wasm", "instantiate_contract", "execute_contract", "migrate_contract", "update_admin", "clear_admin", "multisend", "ibc_transfer":
	// Liveness events
	case "liveness":
		// Not logging these because liveness events are emitted for every block to track validator uptime/missed blocks.
		// High volume, generally only useful for low-level validator monitoring.
	//Mint module events
	case "coinbase", "mint":
		// Not logging these events because coinbase and mint are generated frequently during block production,
		// leading to high event volume mostly relevant for inflation tracking or supply adjustments.
		// Typically not needed for higher-level application indexing.
	// Distribution events
	case "commission", "rewards":
		// Not logging commission and rewards events as they occur very frequently for validator payouts and delegator rewards,
		// which can create significant noise and storage overhead in the indexer.
		// Usually these are monitored via specialized tools or modules.
	// Bank module events
	case "coin_spent", "coin_received", "transfer":
		// Not logging bank module events like coin_spent, coin_received, and transfer because they generate a large volume
		// of events for all token movements across accounts, which can overwhelm the indexer.
		// These are better handled by dedicated transaction or balance tracking systems.
	//IBC Messages
	case "recover_client":
	// Core transaction events
	case "message", "tx":
		// Not logging core transaction events such as message and tx due to their high frequency and verbosity.
		// These events are typically processed by transaction-level handlers or explorers rather than the indexer.
	default:
		if _, loaded := warnedUnhandledAbciEventTypes.LoadOrStore(event.Type, struct{}{}); !loaded {
			s.logger.WithField("type", event.Type).Warn("unhandled abci event type; ignoring")
		}
	}
	return nil
}

// convertEventToMap reconstruct a map based on the event's attribute
func convertEventToMap(event abcitypes.Event) (map[string]any, error) {
	result := make(map[string]any)
	for _, attr := range event.Attributes {
		attrValue := strings.Trim(string(attr.Value), `"`)
		if len(attrValue) == 0 {
			continue
		}
		// Skip handling of "msg_index" field
		if attr.Key == "msg_index" {
			continue
		}
		if attr.Key == "mode" {
			continue
		}
		switch attrValue[0] {
		case '{':
			var nest any
			if err := json.Unmarshal([]byte(attr.Value), &nest); err != nil {
				return nil, fmt.Errorf("fail to unmarshal %s to map,err: %w", attrValue, err)
			}
			result[string(attr.Key)] = nest
		case '[':
			var nest []any
			if err := json.Unmarshal([]byte(attr.Value), &nest); err != nil {
				return nil, fmt.Errorf("fail to unmarshal %s to slice,err: %w", attrValue, err)
			}
			result[string(attr.Key)] = nest
		default:
			result[string(attr.Key)] = attrValue
		}
	}
	return result, nil
}

func subscribe(ctx context.Context, client *tmclient.HTTP, query string) (<-chan ctypes.ResultEvent, error) {
	out, err := client.Subscribe(ctx, "", query)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to query,query:%s,err: %w", query, err)
	}
	return out, nil
}

// parseEventToConcreteType decode all the attribute in given abcitype.Event, and convert it to its relevant concreate type
func parseEventToConcreteType[D any, T interface {
	*D
	proto.Message
}](event abcitypes.Event) (D, error) {
	var defaultValue D
	result, err := convertEventToMap(event)
	if err != nil {
		return defaultValue, err
	}
	buf, err := json.Marshal(result)
	if err != nil {
		return defaultValue, err
	}
	if err := jsonpb.Unmarshal(bytes.NewBuffer(buf), T(&defaultValue)); err != nil {
		return defaultValue, err
	}
	return defaultValue, nil
}
