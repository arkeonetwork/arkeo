package indexer

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"

	abcitypes "github.com/tendermint/tendermint/abci/types"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// type attributeProvider interface {
// 	attributes() map[string]string
// }

type attributes func() map[string]string

func wsAttributeSource(src ctypes.ResultEvent) func() map[string]string {
	attribs := make(map[string]string, len(src.Events))
	for k, v := range src.Events {
		if len(v) > 0 {
			key := k
			if sl := strings.Split(k, "."); len(sl) > 1 {
				key = sl[1]
			}
			if _, ok := attribs[key]; ok {
				log.Debugf("key %s already in results with value %s, overwriting with %s", key, attribs[key], v[0])
			}
			attribs[key] = v[0]
		}
		if len(v) > 1 {
			log.Warnf("attrib %s has %d array values: %v", k, len(v), v)
		}
	}
	attribs["eventHeight"] = attribs["height"]
	return func() map[string]string { return attribs }
}

func tmAttributeSource(tx tmtypes.Tx, evt abcitypes.Event, height int64) func() map[string]string {
	attribs := make(map[string]string, 0)
	for _, attr := range evt.Attributes {
		attribs[string(attr.Key)] = string(attr.Value)
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

func (a *IndexerApp) handleValidatorPayoutEvent(evt types.ValidatorPayoutEvent) error {
	log.Infof("receieved validatorPayoutEvent %#v", evt)
	if evt.Paid < 0 {
		return fmt.Errorf("received negative paid amt: %d for tx %s", evt.Paid, evt.TxID)
	}
	if evt.Paid == 0 {
		return nil
	}
	log.Infof("upserting validator payout event for tx %s", evt.TxID)
	if _, err := a.db.UpsertValidatorPayoutEvent(evt); err != nil {
		return errors.Wrapf(err, "error upserting validator payout event")
	}
	return nil
}

func (a *IndexerApp) consumeEvents(clients []*tmclient.HTTP) error {
	// splitting across multiple tendermint clients as websocket allows max of 5 subscriptions per client
	blockEvents := subscribe(clients[0], "tm.event = 'NewBlock'")
	bondProviderEvents := subscribe(clients[0], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'")
	modProviderEvents := subscribe(clients[0], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'")
	openContractEvents := subscribe(clients[1], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")
	closeContractEvents := subscribe(clients[1], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'")
	claimContractIncomeEvents := subscribe(clients[1], "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	log.Infof("beginning realtime event consumption")
	for {
		select {
		case evt := <-blockEvents:
			data, ok := evt.Data.(tmtypes.EventDataNewBlock)
			if !ok {
				log.Errorf("event not block: %T", evt.Data)
				continue
			}
			log := log.WithField("height", strconv.FormatInt(data.Block.Height, 10))
			log.Debugf("received block: %d", data.Block.Height)

			if err := a.handleBlockEvent(data.Block); err != nil {
				log.Errorf("error handling block event %d: %+v", data.Block.Height, err)
			}

			endBlockEvents := data.ResultEndBlock.Events
			log.Debugf("block %d with %d endBlock events", data.Block.Height, len(endBlockEvents))
			for _, evt := range endBlockEvents {
				switch evt.GetType() {
				case "validator_payout":
					validatorPayoutEvent := types.ValidatorPayoutEvent{}
					if err := convertEvent(tmAttributeSource(nil, evt, data.Block.Height), &validatorPayoutEvent); err != nil {
						log.Errorf("error converting validator_payout event: %+v", err)
						break
					}
					if err := a.handleValidatorPayoutEvent(validatorPayoutEvent); err != nil {
						log.Errorf("error handling validator_payout event: %+v", err)
					}
				case "contract_settlement":
					contractSettlementEvent := types.ContractSettlementEvent{}
					if err := convertEvent(tmAttributeSource(nil, evt, data.Block.Height), &contractSettlementEvent); err != nil {
						log.Errorf("error converting contract_settlement event: %+v", err)
						continue
					}
					// TODO
					// if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
					// 	log.Errorf("error handling close_contract contract_settlement event: %+v", err)
					// }
				}
			}
		case evt := <-openContractEvents:
			log.Debugf("received open contract event")
			if err := a.handleOpenContractEvent(evt); err != nil {
				log.Errorf("error handling open_contract event: %+v", err)
			}
		case evt := <-bondProviderEvents:
			log.Debugf("received bond provider event")
			if err := a.handleBondProviderEvent(evt); err != nil {
				log.Errorf("error handling bond_provider event: %+v", err)
			}
		case evt := <-modProviderEvents:
			log.Debugf("received mod provider event")
			if err := a.handleModProviderEvent(evt); err != nil {
				log.Errorf("error handling mod_provider event: %+v", err)
			}
		case evt := <-claimContractIncomeEvents:
			log.Debugf("received claim contract income event")
			if err := a.handleContractSettlementEvent(evt); err != nil {
				log.Errorf("error handling claim_contract_income event: %+v", err)
			}
		case evt := <-closeContractEvents:
			log.Debugf("received close_contract event")
			if err := a.handleCloseContractEvent(evt); err != nil {
				log.Errorf("error handling close_contract event: %+v", err)
			}

			// TODO needed?
			// if err := a.handleContractSettlementEvent(closeContractEvent.ContractSettlementEvent); err != nil {
			// 	log.Errorf("error handling close_contract contract_settlement event: %+v", err)
			// }
		case <-quit:
			log.Infof("received os quit signal")
			return nil
		}
	}
}

func (a *IndexerApp) consumeHistoricalBlock(client *tmclient.HTTP, bheight int64) (result *db.Block, err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var block *ctypes.ResultBlock
	var blockResults *ctypes.ResultBlockResults
	var blockErr, resultsErr error

	go func() {
		defer wg.Done()
		start := time.Now()
		block, blockErr = client.Block(context.Background(), &bheight)
		if time.Since(start) > 500*time.Millisecond {
			log.Warnf("%.3f elapsed reading block %d", time.Since(start).Seconds(), bheight)
		}
	}()

	go func() {
		defer wg.Done()
		start := time.Now()
		blockResults, resultsErr = client.BlockResults(context.Background(), &bheight)
		if time.Since(start) > 500*time.Millisecond {
			log.Warnf("%.3f elapsed reading block results %d", time.Since(start).Seconds(), bheight)
		}
	}()
	wg.Wait()

	if blockErr != nil {
		return nil, errors.Wrapf(blockErr, "error reading block")
	}
	if resultsErr != nil {
		return nil, errors.Wrapf(resultsErr, "error reading block results")
	}

	log := log.WithField("height", strconv.FormatInt(block.Block.Height, 10))
	for _, transaction := range block.Block.Txs {
		txInfo, err := client.Tx(context.Background(), transaction.Hash(), false)
		if err != nil {
			log.Warnf("failed to get transaction data for %s", transaction.Hash())
			continue
		}

		for _, event := range txInfo.TxResult.Events {
			log.Debugf("received %s txevent", event.Type)
			if err := a.handleAbciEvent(event, transaction, block.Block.Height); err != nil {
				log.Errorf("error handling abci event %#v\n%+v", event, err)
			}
		}
	}

	for _, event := range blockResults.EndBlockEvents {
		log.Debugf("received %s endblock event", event.Type)
		if err := a.handleAbciEvent(event, nil, block.Block.Height); err != nil {
			log.Errorf("error handling abci event %#v\n%+v", event, err)
		}
	}

	r := &db.Block{
		Height:    block.Block.Height,
		Hash:      block.Block.Hash().String(),
		BlockTime: block.Block.Time,
	}
	return r, nil
}

func (a *IndexerApp) handleAbciEvent(event abcitypes.Event, transaction tmtypes.Tx, height int64) error {
	var err error
	switch event.Type {
	case "provider_bond":
		bondProviderEvent := types.BondProviderEvent{}
		if err = convertEvent(tmAttributeSource(transaction, event, height), &bondProviderEvent); err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		// TODO
		// if err = a.handleBondProviderEvent(event); err != nil {
		// 	log.Errorf("error handling %s event: %+v", event.Type, err)
		// }
	case "provider_mod":
		modProviderEvent := types.ModProviderEvent{}
		if err = convertEvent(tmAttributeSource(transaction, event, height), &modProviderEvent); err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		// TODO
		// if err = a.handleModProviderEvent(modProviderEvent); err != nil {
		// 	log.Errorf("error handling %s event: %+v", event.Type, err)
		// }
	case "open_contract":
		openContractEvent := types.OpenContractEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &openContractEvent); err != nil {
			log.Errorf("error converting %s event: %+v", event.Type, err)
			break
		}
		// TODO
		// if err = a.handleOpenContractEvent(openContractEvent); err != nil {
		// 	log.Errorf("error handling %s event: %+v", event.Type, err)
		// }
	case "claim_contract_income":
		contractSettlementEvent := types.ContractSettlementEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &contractSettlementEvent); err != nil {
			log.Errorf("error converting claim_contract_income event: %+v", err)
			break
		}
		// TODO
		// if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
		// 	log.Errorf("error handling claim contract income event: %+v", err)
		// }
	case "validator_payout":
		validatorPayoutEvent := types.ValidatorPayoutEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &validatorPayoutEvent); err != nil {
			log.Errorf("error converting validatorPayoutEvent event: %+v", err)
			break
		}
		if err := a.handleValidatorPayoutEvent(validatorPayoutEvent); err != nil {
			log.Errorf("error handling claim contract income event: %+v", err)
		}
	case "contract_settlement":
		contractSettlementEvent := types.ContractSettlementEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &contractSettlementEvent); err != nil {
			log.Errorf("error converting contractSettlementEvent: %+v", err)
			break
		}
		// TODO
		// if err := a.handleContractSettlementEvent(contractSettlementEvent); err != nil {
		// 	log.Errorf("error handling contractSettlementEvent: %+v", err)
		// }
	case "close_contract":
		log.Debugf("received close_contract event")
		closeContractEvent := types.CloseContractEvent{}
		if err := convertEvent(tmAttributeSource(transaction, event, height), &closeContractEvent); err != nil {
			log.Errorf("error converting close_contract event: %+v", err)
			break
		}
		// TODO
		// if err := a.handleCloseContractEvent(closeContractEvent); err != nil {
		// 	log.Errorf("error handling close contract event: %+v", err)
		// }
	default:
		log.Debugf("ignored event %s", event.Type)
	}
	return nil
}

// copy attributes of map given by attributeFunc() to target which must be a pointer (map/slice implicitly ptr)
func convertEvent(attributeFunc attributes, target interface{}) error {
	return mapstructure.WeakDecode(attributeFunc(), target)
}

func subscribe(client *tmclient.HTTP, query string) <-chan ctypes.ResultEvent {
	out, err := client.Subscribe(context.Background(), "", query)
	if err != nil {
		log.Errorf("failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}
	return out
}
