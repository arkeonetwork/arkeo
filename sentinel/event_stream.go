package sentinel

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/arkeonetwork/arkeo/common"

	"github.com/tendermint/tendermint/libs/log"

	arkeoTypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	tmCoreTypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// TODO: if there are multiple of the same type of event, this may be
// problematic, multiple events may get purged into one (not sure)
func convertEvent(etype string, raw map[string][]string) map[string]string {
	newEvt := make(map[string]string, 0)

	for k, v := range raw {
		if strings.HasPrefix(k, etype+".") {
			parts := strings.SplitN(k, ".", 2)
			newEvt[parts[1]] = v[0]
		}
	}

	return newEvt
}

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

	// receive height changes
	newBlockOut := subscribe(client, logger, "tm.event = 'NewBlockHeader'")
	openContractOut := subscribe(client, logger, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")
	closeContractOut := subscribe(client, logger, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'")
	claimContractOut := subscribe(client, logger, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case result := <-newBlockOut:
			p.handleNewBlockHeaderEvent(result)
		case result := <-openContractOut:
			p.handleOpenContractEvent(result)
		case result := <-closeContractOut:
			p.handleCloseContractEvent(result)
		case result := <-claimContractOut:
			p.handleClaimContractIncomeEvent(result)
		case <-quit:
			return
		}
	}
}

// handleClaimContractIncomeEvent
func (p Proxy) handleClaimContractIncomeEvent(result tmCoreTypes.ResultEvent) {
	evt, err := parseClaimContractIncome(convertEvent(arkeoTypes.EventTypeContractSettlement, result.Events))
	if err != nil {
		p.logger.Error("failed to get close contract event", "error", err)
		return
	}
	if !p.isMyPubKey(evt.Contract.Provider) {
		return
	}
	spender := evt.Contract.Delegate
	if spender.IsEmpty() {
		spender = evt.Contract.Client
	}
	newClaim := NewClaim(evt.Contract.Id, spender, evt.Contract.Nonce, "")
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
	evt, err := parseCloseContract(convertEvent(arkeoTypes.EventTypeCloseContract, result.Events))
	if err != nil {
		p.logger.Error("failed to get close contract event", "error", err)
		return
	}
	if !p.isMyPubKey(evt.Contract.Provider) {
		return
	}
	p.MemStore.Put(evt.Contract)
}

func (p Proxy) handleOpenContractEvent(result tmCoreTypes.ResultEvent) {
	evt, err := parseOpenContract(convertEvent(arkeoTypes.EventTypeOpenContract, result.Events))
	if err != nil {
		p.logger.Error("failed to get open contract event", "error", err)
		return
	}
	if !p.isMyPubKey(evt.Contract.Provider) {
		return
	}
	if evt.Contract.Deposit.IsZero() {
		p.logger.Error("contract's deposit is zero")
		return
	}
	p.MemStore.Put(evt.Contract)
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
		if evt.Type == arkeoTypes.EventTypeContractSettlement {
			input := make(map[string]string)
			for _, attr := range evt.Attributes {
				input[string(attr.Key)] = string(attr.Value)
			}
			evt, err := parseClaimContractIncome(input)
			if err != nil {
				p.logger.Error("failed to get close contract event", "error", err)
				continue
			}
			if !p.isMyPubKey(evt.Contract.Provider) {
				continue
			}
			spender := evt.Contract.Delegate
			if spender.IsEmpty() {
				spender = evt.Contract.Client
			}
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
