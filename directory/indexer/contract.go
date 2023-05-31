package indexer

import (
	"fmt"

	"github.com/pkg/errors"

	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (s *Service) handleOpenContractEvent(evt atypes.EventOpenContract) error {
	provider, err := s.db.FindProvider(evt.Provider.String(), evt.Service)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for service %s", evt.Provider.String(), evt.Service)
	}
	
	_, err = s.db.UpsertContract(provider.ID, evt)
	if err != nil {
		return errors.Wrapf(err, "error upserting contract")
	}
	/*
		// not currently using this
		if _, err = s.db.UpsertOpenContractEvent(ent.ID, evt); err != nil {
			return errors.Wrapf(err, "error upserting open contract event")
		}
	*/

	return nil
}

func (s *Service) handleCloseContractEvent(evt atypes.EventCloseContract, height int64) error {
	/*
		if _, err = s.db.UpsertCloseContractEvent(contract.ID, evt); err != nil {
			return errors.Wrapf(err, "error upserting close contract event")
		}
	*/
	if _, err := s.db.CloseContract(evt.ContractId, height); err != nil {
		return errors.Wrapf(err, "error closing contract %d", evt.ContractId)
	}
	return nil
}

func (s *Service) handleContractSettlementEvent(evt atypes.EventSettleContract) error {
	s.logger.WithField("event", Stringfy(evt)).Info("received event settle contract")
	contract, err := s.db.FindContract(evt.ContractId)
	if err != nil {
		return errors.Wrapf(err, "error finding contract provider %s service %s", evt.Provider, evt.Service)
	}
	if contract == nil {
		return fmt.Errorf("no contract found for provider %s:%s delegPub: %s height %d", evt.Provider, evt.Service, evt.Delegate, evt.Height)
	}
	if _, err = s.db.UpsertContractSettlementEvent(evt); err != nil {
		return errors.Wrapf(err, "error upserting contract settlement event")
	}
	return nil
}
