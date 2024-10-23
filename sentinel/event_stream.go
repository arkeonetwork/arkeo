package sentinel

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gogo/protobuf/proto"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	"github.com/cometbft/cometbft/libs/log"

	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	tmCoreTypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
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
		"tm.event = 'NewBlock'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgBondProvider'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgModProvider'",
	)

	dispatchEvents := func(result tmCoreTypes.ResultEvent) {
		switch {
		case strings.Contains(result.Query, "NewBlock"):
			p.handleNewBlockHeaderEvent(result)

		case strings.Contains(result.Query, "MsgOpenContract"):
			p.handleOpenContractEvent(result)

		case strings.Contains(result.Query, "MsgCloseContract"):
			p.handleCloseContractEvent(result)

		case strings.Contains(result.Query, "MsgClaimContractIncome"):
			p.handleContractSettlementEvent(result)

		case strings.Contains(result.Query, "MsgModProvider"):
			p.handleModProviderEvent(result)

		case strings.Contains(result.Query, "MsgBondProvider"):
			p.handleBondProviderEvent(result)

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
	data, ok := result.Data.(tmtypes.EventDataNewBlock)
	if !ok {
		p.logger.Error("failed cast data")
		return
	}
	height := data.Block.Header.Height
	p.logger.Info("New height detected", "height", height)
	p.MemStore.SetHeight(height)

	for _, evt := range data.ResultFinalizeBlock.Events {
		if evt.Type == types.EventTypeSettleContract {
			input := make(map[string]string)
			for _, attr := range evt.Attributes {
				input[attr.Key] = strings.Trim(attr.Value, `"`)
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

func (p Proxy) handleBondProviderEvent(result tmCoreTypes.ResultEvent) {
	typedEvent, err := parseTypedEvent(result, "arkeo.arkeo.EventBondProvider")
	if err != nil {
		p.logger.Error("failed to parse typed event", "error", err)
		return
	}

	evt, ok := typedEvent.(*types.EventBondProvider)
	if !ok {
		p.logger.Error(fmt.Sprintf("failed to cast %T to EventOpenContract", typedEvent))
		return
	}

	service := common.Service(common.ServiceLookup[evt.Service])
	if !p.isMyPubKey(evt.Provider) {
		return
	}
	providerConfig, err := p.ProviderConfigStore.Get(evt.Provider)
	if err != nil {
		p.logger.Info("failed to get provider config, initializing new config", "error", err)
		providerConfig = ProviderConfiguration{
			PubKey:              evt.Provider,
			Service:             service,
			Bond:                evt.BondAbs,
			BondRelative:        evt.BondRel,
			MetadataUri:         "",
			MetadataNonce:       0,
			Status:              types.ProviderStatus(0),
			MinContractDuration: 0,
			MaxContractDuration: 0,
			SubscriptionRate:    cosmos.Coins{},
			PayAsYouGoRate:      cosmos.Coins{},
			SettlementDuration:  0,
		}
	}

	providerConfig.Bond = evt.BondAbs
	providerConfig.BondRelative = evt.BondRel
	err = p.ProviderConfigStore.Set(providerConfig)
	if err != nil {
		p.logger.Error("failed to update provider configuration", "error", err)
		return
	}
	p.logger.Info("Provider configuration updated on bond provider event", "pubkey", evt.Provider.String(), "service", service.String())

}
func (p Proxy) handleModProviderEvent(result tmCoreTypes.ResultEvent) {
	typedEvent, err := parseTypedEvent(result, "arkeo.arkeo.EventModProvider")
	if err != nil {
		p.logger.Error("failed to parse typed event", "error", err)
		return
	}

	evt, ok := typedEvent.(*types.EventModProvider)
	if !ok {
		p.logger.Error(fmt.Sprintf("failed to cast %T to EventOpenContract", typedEvent))
		return
	}

	service := common.Service(common.ServiceLookup[evt.Service])

	if !p.isMyPubKey(evt.Provider) {
		return
	}

	providerConfig, err := p.ProviderConfigStore.Get(evt.Provider)
	if err != nil {
		p.logger.Error(fmt.Sprintf("failed to get provider %s", err))
		return
	}

	providerConfig.Bond = evt.Bond
	providerConfig.Service = service
	providerConfig.MetadataUri = evt.MetadataUri
	providerConfig.MetadataNonce = evt.MetadataNonce
	providerConfig.Status = evt.Status
	providerConfig.MinContractDuration = evt.MaxContractDuration
	providerConfig.MaxContractDuration = evt.MaxContractDuration
	providerConfig.SubscriptionRate = evt.SubscriptionRate
	providerConfig.PayAsYouGoRate = evt.PayAsYouGoRate
	providerConfig.SettlementDuration = evt.SettlementDuration

	err = p.ProviderConfigStore.Set(providerConfig)
	if err != nil {
		p.logger.Error("failed to update provider configuration", "error", err)
		return
	}
	p.logger.Info("Provider configuration updated on mod provider event", "pubkey", evt.Provider.String(), "service", service.String())

}
