package indexer

import (
	"fmt"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/pkg/errors"
)

func (a *IndexerApp) handleOpenContractEvent(evt types.OpenContractEvent) error {
	provider, err := a.db.FindProvider(evt.ProviderPubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.ProviderPubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("no provider found: DNE %s %s", evt.ProviderPubkey, evt.Chain)
	}
	ent, err := a.db.UpsertContract(provider.ID, evt)
	if err != nil {
		return errors.Wrapf(err, "error upserting contract")
	}
	if _, err = a.db.UpsertOpenContractEvent(ent.ID, evt); err != nil {
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
