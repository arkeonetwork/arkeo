package sentinel

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"cosmossdk.io/errors"
	"github.com/gogo/protobuf/proto"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"

	"github.com/cometbft/cometbft/libs/log"

	tmlog "github.com/cometbft/cometbft/libs/log"
	tmclient "github.com/cometbft/cometbft/rpc/client/http"
	tmCoreTypes "github.com/cometbft/cometbft/rpc/core/types"
	tmtypes "github.com/cometbft/cometbft/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

var numOfWebSocketClients = 2

func subscribe(client *tmclient.HTTP, logger log.Logger, query string) <-chan tmCoreTypes.ResultEvent {
	out, err := client.Subscribe(context.Background(), "", query)
	if err != nil {
		logger.Error("Failed to subscribe to query", "err", err, "query", query)
		os.Exit(1)
	}
	return out
}

func NewTendermintClient(baseURL string, authManager *ArkeoAuthManager) (*tmclient.HTTP, error) {
	// Add auth to WebSocket URL if configured
	if authManager != nil {
		authHeader, err := authManager.GenerateAuthHeader()
		if err != nil {
			return nil, fmt.Errorf("failed to generate auth header: %w", err)
		}

		// Parse URL to add query parameter
		u, err := url.Parse(baseURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse base URL: %w", err)
		}

		q := u.Query()
		q.Set(QueryArkAuth, authHeader)
		u.RawQuery = q.Encode()
		baseURL = u.String()
	}

	client, err := tmclient.New(baseURL, "/websocket")
	if err != nil {
		return nil, errors.Wrapf(err, "error creating websocket client")
	}
	logger := tmlog.NewTMLogger(tmlog.NewSyncWriter(os.Stdout))
	client.SetLogger(logger)

	return client, nil
}

func (p Proxy) EventListener(host string, authManager *ArkeoAuthManager) {
	logger := p.logger

	logger.Info("starting realtime indexing using /websocket")

	// as maximum allowed connection is 5 per ws client(cometbft) we split this into 2 client to handle 3 connection each
	clients := make([]*tmclient.HTTP, numOfWebSocketClients)

	for i := 0; i < numOfWebSocketClients; i++ {
		client, err := NewTendermintClient(fmt.Sprintf("tcp://%s", host), authManager)
		if err != nil {
			panic(fmt.Sprintf("error creating tm client for %s: %+v", host, err))
		}
		if err = client.Start(); err != nil {
			panic(fmt.Sprintf("error starting ws client: %s: %+v", host, err))
		}
		defer func() {
			if err := client.Stop(); err != nil {
				logger.Error("Failed to stop the client", "error", err)
			}
		}()
		clients[i] = client
	}

	// Create a unified channel for receiving events
	eventChan := make(chan tmCoreTypes.ResultEvent, 1000)

	// Function to subscribe to events for a given client
	subscribeToEvents := func(client *tmclient.HTTP, queries ...string) {
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

	// Subscribe to events for each client
	go subscribeToEvents(clients[0],
		"tm.event = 'NewBlock'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'",
		"tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'",
	)

	go subscribeToEvents(clients[1],
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
	providerConfig, err := p.ProviderConfigStore.Get(evt.Provider, service.String())
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

	p.logger.Error(fmt.Sprintf("evt.Provider: ", evt.Provider))
	p.logger.Error(fmt.Sprintf("service.String(): ", service.String()))

	providerConfig, err := p.ProviderConfigStore.Get(evt.Provider, service.String())
	if err != nil {
		p.logger.Error(fmt.Sprintf("failed to get provider %s", err))
		return
	}

	providerConfig.Bond = evt.Bond
	providerConfig.Service = service
	providerConfig.MetadataUri = evt.MetadataUri
	providerConfig.MetadataNonce = evt.MetadataNonce
	providerConfig.Status = evt.Status
	providerConfig.MinContractDuration = evt.MinContractDuration
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
