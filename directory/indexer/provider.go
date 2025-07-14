package indexer

import (
	"context"
	"fmt"
	"github.com/arkeonetwork/arkeo/directory/types"
	"net/url"

	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/utils"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func (s *Service) handleModProviderEvent(ctx context.Context, evt atypes.EventModProvider, txID string, height int64) error {
	provider, err := s.db.FindProvider(ctx, evt.Provider.String(), evt.Service)
	if err != nil {
		return fmt.Errorf("fail to find provider %s for service %s,err: %w", evt.Provider, evt.Service, err)
	}

	log := s.logger.WithField("provider", provider.ID)

	isMetaDataUpdated := provider.MetadataNonce == 0 || provider.MetadataNonce < evt.MetadataNonce
	provider.MetadataURI = evt.MetadataUri
	provider.MetadataNonce = evt.MetadataNonce
	provider.Status = evt.Status.String()
	provider.MinContractDuration = evt.MinContractDuration
	provider.MaxContractDuration = evt.MaxContractDuration
	provider.SubscriptionRate = evt.SubscriptionRate
	provider.PayAsYouGoRate = evt.PayAsYouGoRate
	provider.SettlementDuration = evt.SettlementDuration

	if _, err = s.db.UpdateProvider(ctx, provider); err != nil {
		return fmt.Errorf("error updating provider for mod event %s service %s,err: %w", provider.Pubkey, provider.Service, err)
	}

	// we should create s new modevent struct that looks something
	log.Infof("updated provider %s service %s", provider.Pubkey, provider.Service)
	modEvent := types.ModProviderEvent{
		Pubkey:              evt.Provider.String(),
		Service:             evt.Service,
		Height:              height,
		TxID:                txID,
		MetadataURI:         evt.MetadataUri,
		MetadataNonce:       evt.MetadataNonce,
		Status:              types.ProviderStatus(evt.Status.String()),
		MinContractDuration: evt.MinContractDuration,
		MaxContractDuration: evt.MaxContractDuration,
		SettlementDuration:  evt.SettlementDuration,
		SubscriptionRate:    evt.SubscriptionRate,
		PayAsYouGoRate:      evt.PayAsYouGoRate,
	}
	if _, err = s.db.InsertModProviderEvent(ctx, provider.ID, modEvent, txID, height); err != nil {
		return errors.Wrapf(err, "error inserting ModProviderEvent for %s service %s", evt.Provider, evt.Service)
	}

	if !isMetaDataUpdated {
		return nil
	}

	log.Debugf("updating provider metadata for provider %s", provider.Pubkey)
	if !validateMetadataURI(provider.MetadataURI) {
		log.Errorf("updating provider metadata for provider %s failed due to bad MetadataURI %s", provider.Pubkey, provider.MetadataURI)
		return nil
	}
	providerMetadata, err := utils.DownloadProviderMetadata(provider.MetadataURI, 5, 1e6)
	if err != nil {
		log.WithError(err).Errorf("updating provider metadata for provider %s failed", provider.Pubkey)
		return nil
	}

	if providerMetadata == nil {
		log.Errorf("nil providerMetadata for %s", provider.MetadataURI)
		return nil
	}

	if _, err = s.db.UpsertProviderMetadata(ctx, provider.ID, int64(provider.MetadataNonce), *providerMetadata); err != nil {
		return errors.Wrapf(err, "error updating provider metadta for mod event %s service %s", provider.Pubkey, provider.Service)
	}

	return nil
}

func (s *Service) handleBondProviderEvent(ctx context.Context, evt atypes.EventBondProvider, txID string, height int64) error {
	isNewProvider := false
	provider, err := s.db.FindProvider(ctx, evt.Provider.String(), evt.Service)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			return errors.Wrapf(err, "error finding provider %s for service %s", evt.Provider, evt.Service)
		}

		// provider doesn't exist yet , create a new one
		provider, err = s.createProvider(ctx, evt)
		if err != nil {
			return errors.Wrapf(err, "error creating provider %s service %s", evt.Provider, evt.Service)
		}
		isNewProvider = true
	}
	if !isNewProvider {
		if evt.BondAbs.IsNil() {
			provider.Bond = evt.BondAbs.String()
		}
		// TODO change this to just update bond , `UpdateProvider` does a lot other stuff
		if _, err = s.db.UpdateProvider(ctx, provider); err != nil {
			return errors.Wrapf(err, "error updating provider for bond event %s service %s", evt.Provider, evt.Service)
		}
	}

	s.logger.Debugf("handled bond provider event for %s service %s", evt.Provider, evt.Service)
	if _, err = s.db.InsertBondProviderEvent(ctx, provider.ID, evt, height, txID); err != nil {
		return errors.Wrapf(err, "error inserting BondProviderEvent for %s service %s", evt.Provider, evt.Service)
	}
	return nil
}

func (s *Service) createProvider(ctx context.Context, evt atypes.EventBondProvider) (*db.ArkeoProvider, error) {
	// new provider for service, insert
	provider := &db.ArkeoProvider{
		Pubkey:  evt.Provider.String(),
		Service: evt.Service,
		Bond:    evt.BondAbs.String(),
	}
	entity, err := s.db.InsertProvider(ctx, provider)
	if err != nil {
		return nil, fmt.Errorf("fail to insert provider %s %s,err: %w", evt.Provider, evt.Service, err)
	}
	if entity == nil {
		return nil, fmt.Errorf("nil entity after inserting provider")
	}
	s.logger.Debugf("inserted provider record %d for %s %s", entity.ID, evt.Provider, evt.Service)
	provider.Entity = *entity
	return provider, nil
}

func validateMetadataURI(uri string) bool {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return false
	}
	return true
}
