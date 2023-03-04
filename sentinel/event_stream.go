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

	tmclient "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
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

func subscribe(client *tmclient.HTTP, logger log.Logger, query string) <-chan ctypes.ResultEvent {
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
		logger.Error("failure to create websocket cliennt", "error", err)
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
	heightOut := subscribe(client, logger, "tm.event = 'NewBlockHeader'")
	openContractOut := subscribe(client, logger, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgOpenContract'")
	closeContractOut := subscribe(client, logger, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgCloseContract'")
	claimContractOut := subscribe(client, logger, "tm.event = 'Tx' AND message.action='/arkeo.arkeo.MsgClaimContractIncome'")
	// miscOut := subscribe(client, logger, "tm.event = 'Tx'")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	isMyPubKey := func(pk common.PubKey) bool {
		return pk.Equals(p.Config.ProviderPubKey)
	}

	for {
		select {
		/*
			case result := <-miscOut:
				for k, v := range result.Events {
					// fmt.Printf("Tx: %s --> %s\n", k, v[0])
				}
		*/
		case result := <-heightOut:
			data, ok := result.Data.(tmtypes.EventDataNewBlockHeader)
			if !ok {
				logger.Error("failed cast data")
				continue
			}
			height := data.Header.Height
			logger.Info("New height detected", "height", height)
			p.MemStore.SetHeight(height)

			// for _, evt := range data.ResultBeginBlock.Events {}

			for _, evt := range data.ResultEndBlock.Events {
				if evt.Type == "contract_settlement" {
					input := make(map[string]string)
					for _, attr := range evt.Attributes {
						input[string(attr.Key)] = string(attr.Value)
					}
					evt, err := parseClaimContractIncome(input)
					if err != nil {
						logger.Error("failed to get close contract event", "error", err)
						continue
					}
					if !isMyPubKey(evt.Contract.ProviderPubKey) {
						continue
					}
					spender := evt.Contract.Delegate
					if spender.IsEmpty() {
						spender = evt.Contract.Client
					}
					newClaim := NewClaim(evt.Contract.ProviderPubKey, evt.Contract.Id, spender, evt.Contract.Nonce, evt.Contract.Height, "")
					currClaim, err := p.ClaimStore.Get(newClaim.Key())
					if err != nil {
						logger.Error("failed to get claim", "error", err)
						continue
					}
					if currClaim.Nonce == newClaim.Nonce && currClaim.Height == newClaim.Height {
						currClaim.Claimed = true
						if err := p.ClaimStore.Set(currClaim); err != nil {
							logger.Error("failed to set claimed", "error", err)
						}
					}
				}
			}
		case result := <-openContractOut:
			evt, err := parseOpenContract(convertEvent("open_contract", result.Events))
			if err != nil {
				logger.Error("failed to get open contract event", "error", err)
				continue
			}
			if !isMyPubKey(evt.Contract.ProviderPubKey) {
				continue
			}
			if evt.Contract.Deposit.IsZero() {
				logger.Error("contract's deposit is zero")
				continue
			}
			p.MemStore.Put(evt.Contract)
		case result := <-closeContractOut:
			evt, err := parseCloseContract(convertEvent("close_contract", result.Events))
			if err != nil {
				logger.Error("failed to get close contract event", "error", err)
				continue
			}
			if !isMyPubKey(evt.Contract.ProviderPubKey) {
				continue
			}
			p.MemStore.Put(evt.Contract)
		case result := <-claimContractOut:
			evt, err := parseClaimContractIncome(convertEvent("contract_settlement", result.Events))
			if err != nil {
				logger.Error("failed to get close contract event", "error", err)
				continue
			}
			if !isMyPubKey(evt.Contract.ProviderPubKey) {
				continue
			}
			spender := evt.Contract.Delegate
			if spender.IsEmpty() {
				spender = evt.Contract.Client
			}
			newClaim := NewClaim(evt.Contract.ProviderPubKey, evt.Contract.Id, spender, evt.Contract.Nonce, evt.Contract.Height, "")
			currClaim, err := p.ClaimStore.Get(newClaim.Key())
			if err != nil {
				logger.Error("failed to get claim", "error", err)
				continue
			}
			if currClaim.Nonce == newClaim.Nonce && currClaim.Height == newClaim.Height {
				currClaim.Claimed = true
				if err := p.ClaimStore.Set(currClaim); err != nil {
					logger.Error("failed to set claimed", "error", err)
				}
			}
		case <-quit:
			break
		}
	}
}
