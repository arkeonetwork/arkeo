package indexer

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/directory/utils"
	atypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/pkg/errors"
)

func (s *Service) handleModProviderEvent(evt atypes.EventModProvider) error {
	provider, err := s.db.FindProvider(evt.Provider.String(), evt.Service)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for service %s", evt.Provider, evt.Service)
	}
	if provider == nil {
		return fmt.Errorf("cannot mod provider, DNE %s %s", evt.Provider, evt.Service)
	}

	log := s.logger.WithField("provider", strconv.FormatInt(provider.ID, 10))

	isMetaDataUpdated := provider.MetadataNonce == 0 || provider.MetadataNonce < evt.MetadataNonce
	provider.MetadataURI = evt.MetadataUri
	provider.MetadataNonce = evt.MetadataNonce
	provider.Status = types.ProviderStatus(evt.Status.String())
	provider.MinContractDuration = evt.MinContractDuration
	provider.MaxContractDuration = evt.MaxContractDuration
	provider.SubscriptionRate = evt.SubscriptionRate
	provider.PayAsYouGoRate = evt.PayAsYouGoRate
	provider.SettlementDuration = evt.SettlementDuration

	if _, err = s.db.UpdateProvider(provider); err != nil {
		return errors.Wrapf(err, "error updating provider for mod event %s service %s", provider.Pubkey, provider.Service)
	}
	/*
		// currently, we're not utilizing the inserts for mod provider events, so
		// i'm disabling for now. If we want to re-enable it (because we see s need
		// for it), we should create s new modevent struct that looks something
		// like this...
		type ModProviderEvent struct {
			atypes.EventModProvider
			Height              int64          `mapstructure:"height"`
			TxID                string         `mapstructure:"hash"`
		}

		log.Infof("updated provider %s service %s", provider.Pubkey, provider.Service)
		if _, err = s.db.InsertModProviderEvent(provider.ID, evt); err != nil {
			return errors.Wrapf(err, "error inserting ModProviderEvent for %s service %s", evt.Provider, evt.Service)
		}
	*/

	if !isMetaDataUpdated {
		return nil
	}

	log.Debugf("updating provider metadata for provider %s", provider.Pubkey)
	if !validateMetadataURI(provider.MetadataURI) {
		log.Warnf("updating provider metadata for provider %s failed due to bad MetadataURI %s", provider.Pubkey, provider.MetadataURI)
		return nil
	}
	providerMetadata, err := utils.DownloadProviderMetadata(provider.MetadataURI, 5, 1e6)
	if err != nil {
		log.Warnf("updating provider metadata for provider %s failed %v", provider.Pubkey, err)
		return nil
	}

	if providerMetadata == nil {
		log.Errorf("nil providerMetadata for %s", provider.MetadataURI)
		return nil
	}

	if _, err = s.db.UpsertProviderMetadata(provider.ID, int64(provider.MetadataNonce), *providerMetadata); err != nil {
		return errors.Wrapf(err, "error updating provider metadta for mod event %s service %s", provider.Pubkey, provider.Service)
	}
	return nil
}

func (s *Service) handleBondProviderEvent(evt types.BondProviderEvent) error {
	provider, err := s.db.FindProvider(evt.Pubkey, evt.Service)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for service %s", evt.Pubkey, evt.Service)
	}
	if provider == nil {
		// new provider for service, insert
		if provider, err = s.createProvider(evt); err != nil {
			return errors.Wrapf(err, "error creating provider %s service %s", evt.Pubkey, evt.Service)
		}
	} else {
		if evt.BondAbsolute != "" {
			provider.Bond = evt.BondAbsolute
		}
		if _, err = s.db.UpdateProvider(provider); err != nil {
			return errors.Wrapf(err, "error updating provider for bond event %s service %s", evt.Pubkey, evt.Service)
		}
	}

	s.logger.Debugf("handled bond provider event for %s service %s", evt.Pubkey, evt.Service)
	if _, err = s.db.InsertBondProviderEvent(provider.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting BondProviderEvent for %s service %s", evt.Pubkey, evt.Service)
	}
	return nil
}

func (s *Service) createProvider(evt types.BondProviderEvent) (*db.ArkeoProvider, error) {
	// new provider for service, insert
	provider := &db.ArkeoProvider{Pubkey: evt.Pubkey, Service: evt.Service, Bond: evt.BondAbsolute}
	entity, err := s.db.InsertProvider(provider)
	if err != nil {
		return nil, errors.Wrapf(err, "error inserting provider %s %s", evt.Pubkey, evt.Service)
	}
	if entity == nil {
		return nil, fmt.Errorf("nil entity after inserting provider")
	}
	s.logger.Debugf("inserted provider record %d for %s %s", entity.ID, evt.Pubkey, evt.Service)
	provider.Entity = *entity
	return provider, nil
}

func validateMetadataURI(uri string) bool {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return false
	}
	return true
}
