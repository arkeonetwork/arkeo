package indexer

import (
	"encoding/hex"
	"fmt"
	"strings"

	arkeoUtils "github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/types"
	arkeoTypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (a *IndexerApp) handleOpenContractEvent(evt ctypes.ResultEvent) error {
	typedEvent, err := arkeoUtils.ParseTypedEvent(evt, "arkeo.arkeo.EventOpenContract")
	if err != nil {
		log.Errorf("failed to parse typed event", "error", err)
		return errors.Wrapf(err, "failed to parse typed event")
	}

	txData, ok := evt.Data.(tmtypes.EventDataTx)
	if !ok {
		return fmt.Errorf("failed to cast %T to EventDataTx", evt.Data)
	}

	txid := strings.ToUpper(hex.EncodeToString(tmtypes.Tx(txData.Tx).Hash()))

	openContractEvent, ok := typedEvent.(*arkeoTypes.EventOpenContract)
	if !ok {
		return fmt.Errorf("failed to cast %T to EventOpenContract", typedEvent)
	}

	provider, err := a.db.FindProvider(openContractEvent.Provider.String(), openContractEvent.Service)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", openContractEvent.Provider.String(), openContractEvent.Service)
	}
	if provider == nil {
		return fmt.Errorf("no provider found: DNE %s %s", openContractEvent.Provider.String(), openContractEvent.Service)
	}

	height := openContractEvent.Height

	oce := types.OpenContractEvent{
		BaseContractEvent: types.BaseContractEvent{
			ProviderPubkey: openContractEvent.Provider.String(),
			Chain:          openContractEvent.Service,
			ClientPubkey:   openContractEvent.Client.String(),
			DelegatePubkey: openContractEvent.Delegate.String(),
			TxID:           txid,
			Height:         height,
			EventHeight:    height,
		},
		ContractType: types.ContractType(openContractEvent.Type.String()),
		Duration:     openContractEvent.Duration,
		Rate:         openContractEvent.Rate.Amount.Int64(),
		OpenCost:     openContractEvent.OpenCost,
	}
	ent, err := a.db.UpsertContract(provider.ID, oce)
	if err != nil {
		return errors.Wrapf(err, "error upserting contract")
	}
	if _, err = a.db.UpsertOpenContractEvent(ent.ID, oce); err != nil {
		return errors.Wrapf(err, "error upserting open contract event")
	}

	return nil
}

func (a *IndexerApp) handleCloseContractEvent(evt ctypes.ResultEvent) error {

	typedEvent, err := arkeoUtils.ParseTypedEvent(evt, "arkeo.arkeo.EventCloseContract")
	if err != nil {
		log.Errorf("failed to parse typed event", "error", err)
		return errors.Wrapf(err, "failed to parse typed event")
	}

	txData, ok := evt.Data.(tmtypes.EventDataTx)
	if !ok {
		return fmt.Errorf("failed to cast %T to EventDataTx", evt.Data)
	}

	txid := strings.ToUpper(hex.EncodeToString(tmtypes.Tx(txData.Tx).Hash()))

	closeContractEvent, ok := typedEvent.(*arkeoTypes.EventCloseContract)
	if !ok {
		return fmt.Errorf("failed to cast %T to EventCloseContract", typedEvent)
	}

	clientPubkey := closeContractEvent.Client.String()
	if closeContractEvent.Delegate != "" {
		clientPubkey = closeContractEvent.Delegate.String()
	}

	// FindContractsByPubKeys returns by id descending (newest)
	contracts, err := a.db.FindContractsByPubKeys(closeContractEvent.Service, closeContractEvent.Provider.String(), clientPubkey)
	if err != nil {
		return errors.Wrapf(err, "error finding contract for %s:%s %s", closeContractEvent.Provider, closeContractEvent.Service, clientPubkey)
	}
	if len(contracts) < 1 {
		return fmt.Errorf("no contracts found: %s:%s %s", closeContractEvent.Provider, closeContractEvent.Service, clientPubkey)
	}

	contract := contracts[0]
	event := types.CloseContractEvent{
		ContractSettlementEvent: types.ContractSettlementEvent{
			BaseContractEvent: types.BaseContractEvent{
				ProviderPubkey: closeContractEvent.Provider.String(),
				Chain:          closeContractEvent.Service,
				ClientPubkey:   closeContractEvent.Client.String(),
				DelegatePubkey: closeContractEvent.Delegate.String(),
				TxID:           txid,
				Height:         contract.Height,
				EventHeight:    a.Height,
			},
		},
	}
	if _, err = a.db.UpsertCloseContractEvent(contract.ID, event); err != nil {
		return errors.Wrapf(err, "error upserting open contract event")
	}

	if _, err = a.db.CloseContract(contract.ID, event.EventHeight); err != nil {
		return errors.Wrapf(err, "error closing contract %d", contract.ID)
	}
	return nil
}

func (a *IndexerApp) handleContractSettlementEvent(evt ctypes.ResultEvent) error {
	log.Infof("receieved contractSettlementEvent %#v", evt)

	typedEvent, err := arkeoUtils.ParseTypedEvent(evt, "arkeo.arkeo.EventSettleContract")
	if err != nil {
		log.Errorf("failed to parse typed event: %+v", err)
		return errors.Wrapf(err, "failed to parse typed event")
	}

	txData, ok := evt.Data.(tmtypes.EventDataTx)
	if !ok {
		return fmt.Errorf("failed to cast %T to EventDataTx", evt.Data)
	}

	txid := strings.ToUpper(hex.EncodeToString(tmtypes.Tx(txData.Tx).Hash()))

	settleContractEvent, ok := typedEvent.(*arkeoTypes.EventSettleContract)
	if !ok {
		return fmt.Errorf("failed to cast %T to EventSettleContract", typedEvent)
	}
	clientPubkey := settleContractEvent.Client.String()
	if settleContractEvent.Delegate != "" {
		clientPubkey = settleContractEvent.Delegate.String()
	}

	contract, err := a.db.FindContractByPubKeys(settleContractEvent.Service, settleContractEvent.Provider.String(), clientPubkey, settleContractEvent.Height)
	if err != nil {
		return errors.Wrapf(err, "error finding contract provider %s chain %s", settleContractEvent.Provider, settleContractEvent.Service)
	}
	if contract == nil {
		return fmt.Errorf("no contract found for provider %s:%s delegPub: %s height %d", settleContractEvent.Provider, settleContractEvent.Service, clientPubkey, settleContractEvent.Height)
	}

	event := types.ContractSettlementEvent{
		BaseContractEvent: types.BaseContractEvent{
			ProviderPubkey: settleContractEvent.Provider.String(),
			Chain:          settleContractEvent.Service,
			ClientPubkey:   settleContractEvent.Client.String(),
			DelegatePubkey: settleContractEvent.Delegate.String(),
			TxID:           txid,
			Height:         contract.Height,
			EventHeight:    settleContractEvent.Height,
		},
		Nonce:   fmt.Sprintf("%d", settleContractEvent.Nonce),
		Paid:    settleContractEvent.Paid.String(),
		Reserve: settleContractEvent.Reserve.String(),
	}

	if _, err = a.db.UpsertContractSettlementEvent(contract.ID, event); err != nil {
		return errors.Wrapf(err, "error upserting contract settlement event")
	}
	return nil
}
