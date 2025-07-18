package indexer

import (
	"context"

	"github.com/pkg/errors"
)

func (s *Service) handleGenericEvent(ctx context.Context, eventType string, txID string, height int64, attrJSON []byte) error {
	if _, err := s.db.InsertGenericEvent(ctx, eventType, txID, height, attrJSON); err != nil {
		return errors.Wrapf(err, "error inserting generic event")
	}
	return nil
}
