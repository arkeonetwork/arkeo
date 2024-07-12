package sentinel

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/gogo/protobuf/proto"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func subscribe(client *tmclient.HTTP, logger log.Logger, query string) <-chan tmCoreTypes.ResultEvent {
	out, err := client.Subscribe(context.Background(), "", query)
	if err != nil {
		logger.Error("Failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}
	return out
}

func (p Proxy) EventListener(host string) {
	logger := p.logger
	client, err := tmclient.New(fmt.Sprintf("tcp://%s", host), "/websocket")
	if err != nil {
		logger.Error("failure to create websocket client", "error", err)
		panic(err)
	}
	client.SetLogger(logger)
	err = client.Start()
	if err != nil {
		logger.Error("Failed to start a client", "err", err)
		os.Exit(1)
	}
	defer client.Stop() // nolint

	// create a unified channel for receiving events

	eventChan := make(chan tmCoreTypes.ResultEvent, 1000)

	subscribeToEvents := func(queries ...string) {
		for _, query := range queries {
			out := subscribe(client, logger, query)

			go func(out <-chan tmCoreTypes.ResultEvent) {
				for {
					select {
					case result := <-out:
						eventChan <- result

					case <-client.Quit():
						return
					}
				}
			}(out)
		}
	}

	// subscribe to events
	go subscribeToEvents(
		"tm.event = 'NewBlockHeader'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'",
	)

	dispatchEvents := func(result tmCoreTypes.ResultEvent) {
		switch {
		case strings.Contains(result.Query, "NewBlockHeader"):
			p.handleNewBlockHeaderEvent(result)

		case strings.Contains(result.Query, "MsgOpenContract"):
			p.handleOpenContractEvent(result)

		case strings.Contains(result.Query, "MsgCloseContract"):
			p.handleCloseContractEvent(result)

		case strings.Contains(result.Query, "MsgClaimContractIncome"):
			p.handleContractSettlementEvent(result)

		default:
			logger.Error("Unknown Event Type", "Query", result.Query)
		}
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case result := <-eventChan:
			dispatchEvents(result)
		case <-quit:
			return
		}
	}
}

// handleContractSettlementEvent
func (p Proxy) handleContractSettlementEvent(result tmCoreTypes.ResultEvent) {
	typedEvent, err := parseTypedEvent(result, "arkeo.arkeo.EventSettleContract")
	if err != nil {
		p.logger.Error("failed to parse typed event", "error", err)
		return
	}

	evt, ok := typedEvent.(*types.EventSettleContract)
	if !ok {
		p.logger.Error(fmt.Sprintf("failed to cast %T to EventSettleContract", typedEvent))
		return
	}

	if !p.isMyPubKey(evt.Provider) {
		return
	}

	service := common.Service(common.ServiceLookup[evt.Service])
	contract := types.Contract{
		Provider: evt.Provider,
		Service:  service,
		Client:   evt.Client,
		Delegate: evt.Delegate,
		Id:       evt.ContractId,
	}

	spender := contract.GetSpender()
	newClaim := NewClaim(contract.Id, spender, evt.Nonce, "")
	currClaim, err := p.ClaimStore.Get(newClaim.Key())
	if err != nil {
		p.logger.Error("failed to get claim", "error", err)
		return
	}
	if currClaim.Nonce == newClaim.Nonce {
		currClaim.Claimed = true
		if err := p.ClaimStore.Set(currClaim); err != nil {
			p.logger.Error("failed to set claimed", "error", err)
		}
	}
}

// handleCloseContractEvent
func (p Proxy) handleCloseContractEvent(result tmCoreTypes.ResultEvent) {
	typedEvent, err := parseTypedEvent(result, "arkeo.arkeo.EventCloseContract")
	if err != nil {
		p.logger.Error("failed to parse typed event", "error", err)
		return
	}

	evt, ok := typedEvent.(*types.EventCloseContract)
	if !ok {
		p.logger.Error(fmt.Sprintf("failed to cast %T to EventCloseContract", typedEvent))
		return
	}

	service := common.Service(common.ServiceLookup[evt.Service])
	contract := types.Contract{
		Provider: evt.Provider,
		Service:  service,
		Client:   evt.Client,
		Delegate: evt.Delegate,
		Id:       evt.ContractId,
	}
	if !p.isMyPubKey(contract.Provider) {
		return
	}
	p.MemStore.Put(contract)
}

func (p Proxy) handleOpenContractEvent(result tmCoreTypes.ResultEvent) {
	typedEvent, err := parseTypedEvent(result, "arkeo.arkeo.EventOpenContract")
	if err != nil {
		p.logger.Error("failed to parse typed event", "error", err)
		return
	}

	evt, ok := typedEvent.(*types.EventOpenContract)
	if !ok {
		p.logger.Error(fmt.Sprintf("failed to cast %T to EventOpenContract", typedEvent))
		return
	}

	service := common.Service(common.ServiceLookup[evt.Service])
	contract := types.Contract{
		Provider:           evt.Provider,
		Service:            service,
		Client:             evt.Client,
		Delegate:           evt.Delegate,
		Type:               evt.Type,
		Height:             evt.Height,
		Duration:           evt.Duration,
		Rate:               evt.Rate,
		Deposit:            evt.Deposit,
		Id:                 evt.ContractId,
		SettlementDuration: evt.SettlementDuration,
		Authorization:      evt.Authorization,
		QueriesPerMinute:   evt.QueriesPerMinute,
	}

	if !p.isMyPubKey(evt.Provider) {
		return
	}
	if evt.Deposit.IsZero() {
		p.logger.Error("contract's deposit is zero")
		return
	}
	p.MemStore.Put(contract)
}

func (p Proxy) handleNewBlockHeaderEvent(result tmCoreTypes.ResultEvent) {
	data, ok := result.Data.(tmtypes.EventDataNewBlockHeader)
	if !ok {
		p.logger.Error("failed cast data")
		return
	}
	height := data.Header.Height
	p.logger.Info("New height detected", "height", height)
	p.MemStore.SetHeight(height)

	for _, evt := range data.ResultEndBlock.Events {
		if evt.Type == types.EventTypeSettleContract {
			input := make(map[string]string)
			for _, attr := range evt.Attributes {
				input[string(attr.Key)] = strings.Trim(string(attr.Value), `"`)
			}
			evt, err := parseContractSettlementEvent(input)
			if err != nil {
				p.logger.Error("failed to parse contract settlement event", "error", err)
				continue
			}
			if !p.isMyPubKey(evt.Contract.Provider) {
				continue
			}
			spender := evt.Contract.GetSpender()
			newClaim := NewClaim(evt.Contract.Id, spender, evt.Contract.Nonce, "")
			currClaim, err := p.ClaimStore.Get(newClaim.Key())
			if err != nil {
				p.logger.Error("failed to get claim", "error", err)
				continue
			}
			if currClaim.Nonce == newClaim.Nonce {
				currClaim.Claimed = true
				if err := p.ClaimStore.Set(currClaim); err != nil {
					p.logger.Error("failed to set claimed", "error", err)
				}
			}
		}
	}
}

func (p Proxy) isMyPubKey(pk common.PubKey) bool {
	return pk.Equals(p.Config.ProviderPubKey)
}

func parseTypedEvent(result tmCoreTypes.ResultEvent, eventType string) (proto.Message, error) {
	var (
		msg         proto.Message
		eventDataTx tmtypes.EventDataTx
		ok          bool
	)
	if eventDataTx, ok = result.Data.(tmtypes.EventDataTx); !ok {
		return msg, fmt.Errorf("failed cast %T to EventDataTx", result.Data)
	}

	for _, evt := range eventDataTx.TxResult.Result.Events {
		if evt.Type == eventType {
			return sdk.ParseTypedEvent(evt)
		}
	}

	return msg, fmt.Errorf("event %s not found", eventType)
}
