package indexer

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/directory/types"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (a *IndexerApp) handleOpenContractEvent(evt atypes.EventOpenContract) error {
	provider, err := a.db.FindProvider(evt.Provider.String(), evt.Service)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for service %s", evt.Provider.String(), evt.Service)
	}
	if provider == nil {
		return fmt.Errorf("no provider found: DNE %s %s", evt.Provider.String(), evt.Service)
	}
	_, err = a.db.UpsertContract(provider.ID, evt)
	if err != nil {
		return errors.Wrapf(err, "error upserting contract")
	}
	/*
		// not currently using this
		if _, err = a.db.UpsertOpenContractEvent(ent.ID, evt); err != nil {
			return errors.Wrapf(err, "error upserting open contract event")
		}
	*/

	return nil
}

func (a *IndexerApp) handleCloseContractEvent(evt types.CloseContractEvent) error {
	/*
		if _, err = a.db.UpsertCloseContractEvent(contract.ID, evt); err != nil {
			return errors.Wrapf(err, "error upserting close contract event")
		}
	*/
	if _, err := a.db.CloseContract(evt.ContractId, evt.EventHeight); err != nil {
		return errors.Wrapf(err, "error closing contract %d", evt.ContractId)
	}
	return nil
}

func (a *IndexerApp) handleContractSettlementEvent(evt types.ContractSettlementEvent) error {
	log.Infof("receieved contractSettlementEvent %#v", evt)
	contract, err := a.db.FindContract(evt.ContractId)
	if err != nil {
		return errors.Wrapf(err, "error finding contract provider %s service %s", evt.ProviderPubkey, evt.Service)
	}
	if contract == nil {
		return fmt.Errorf("no contract found for provider %s:%s delegPub: %s height %d", evt.ProviderPubkey, evt.Service, evt.GetDelegatePubkey(), evt.Height)
	}
	if _, err = a.db.UpsertContractSettlementEvent(evt); err != nil {
		return errors.Wrapf(err, "error upserting contract settlement event")
	}
	return nil
}
