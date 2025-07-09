package indexer

import (
	"context"

	"github.com/pkg/errors"

	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (s *Service) handleOpenContractEvent(ctx context.Context, evt atypes.EventOpenContract, txID string, height int64) error {
	provider, err := s.db.FindProvider(ctx, evt.Provider.String(), evt.Service)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for service %s", evt.Provider.String(), evt.Service)
	}

	_, err = s.db.UpsertContract(ctx, provider.ID, evt, txID, height)
	if err != nil {
		return errors.Wrapf(err, "error upserting contract")
	}

	return nil
}

func (s *Service) handleCloseContractEvent(ctx context.Context, evt atypes.EventCloseContract, txID string, height int64) error {
	if _, err := s.db.CloseContract(ctx, evt.ContractId, txID, height); err != nil {
		return errors.Wrapf(err, "error closing contract %d", evt.ContractId)
	}
	return nil
}

func (s *Service) handleContractSettlementEvent(ctx context.Context, evt atypes.EventSettleContract, txID string, height int64) error {
	if _, err := s.db.UpsertContractSettlementEvent(ctx, evt, txID, height); err != nil {
		return errors.Wrapf(err, "error upserting contract settlement event")
	}
	return nil
}
