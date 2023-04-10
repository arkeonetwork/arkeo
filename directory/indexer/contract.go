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

func (a *IndexerApp) handleCloseContractEvent(evt types.CloseContractEvent) error {
	contracts, err := a.db.FindContractsByPubKeys(evt.Chain, evt.ProviderPubkey, evt.GetDelegatePubkey())
	if err != nil {
		return errors.Wrapf(err, "error finding contract for %s:%s %s", evt.ProviderPubkey, evt.Chain, evt.GetDelegatePubkey())
	}
	if len(contracts) < 1 {
		return fmt.Errorf("no contracts found: %s:%s %s", evt.ProviderPubkey, evt.Chain, evt.GetDelegatePubkey())
	}

	// FindContractsByPubKeys returns by id descending (newest)
	contract := contracts[0]
	if _, err = a.db.UpsertCloseContractEvent(contract.ID, evt); err != nil {
		return errors.Wrapf(err, "error upserting open contract event")
	}

	if _, err = a.db.CloseContract(contract.ID, evt.EventHeight); err != nil {
		return errors.Wrapf(err, "error closing contract %d", contract.ID)
	}
	return nil
}

func (a *IndexerApp) handleContractSettlementEvent(evt types.ContractSettlementEvent) error {
	log.Infof("receieved contractSettlementEvent %#v", evt)
	contract, err := a.db.FindContractByPubKeys(evt.Chain, evt.ProviderPubkey, evt.GetDelegatePubkey(), evt.Height)
	if err != nil {
		return errors.Wrapf(err, "error finding contract provider %s chain %s", evt.ProviderPubkey, evt.Chain)
	}
	if contract == nil {
		return fmt.Errorf("no contract found for provider %s:%s delegPub: %s height %d", evt.ProviderPubkey, evt.Chain, evt.GetDelegatePubkey(), evt.Height)
	}
	if _, err = a.db.UpsertContractSettlementEvent(contract.ID, evt); err != nil {
		return errors.Wrapf(err, "error upserting contract settlement event")
	}
	return nil
}
