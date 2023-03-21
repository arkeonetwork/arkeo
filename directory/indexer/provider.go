package indexer

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/directory/utils"
	"github.com/pkg/errors"
)

func (a *IndexerApp) handleModProviderEvent(evt types.ModProviderEvent) error {
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		return fmt.Errorf("cannot mod provider, DNE %s %s", evt.Pubkey, evt.Chain)
	}

	log := log.WithField("provider", strconv.FormatInt(provider.ID, 10))

	isMetaDataUpdated := provider.MetadataNonce == 0 || provider.MetadataNonce < evt.MetadataNonce
	provider.MetadataURI = evt.MetadataURI
	provider.MetadataNonce = evt.MetadataNonce
	provider.Status = evt.Status
	provider.MinContractDuration = evt.MinContractDuration
	provider.MaxContractDuration = evt.MaxContractDuration
	provider.SubscriptionRate = evt.SubscriptionRate
	provider.PayAsYouGoRate = evt.PayAsYouGoRate

	if _, err = a.db.UpdateProvider(provider); err != nil {
		return errors.Wrapf(err, "error updating provider for mod event %s chain %s", provider.Pubkey, provider.Chain)
	}
	log.Infof("updated provider %s chain %s", provider.Pubkey, provider.Chain)
	if _, err = a.db.InsertModProviderEvent(provider.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting ModProviderEvent for %s chain %s", evt.Pubkey, evt.Chain)
	}

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

	providerMetadata.Configuration.Nonce = int64(provider.MetadataNonce)
	if _, err = a.db.UpsertProviderMetadata(provider.ID, *providerMetadata); err != nil {
		return errors.Wrapf(err, "error updating provider metadta for mod event %s chain %s", provider.Pubkey, provider.Chain)
	}
	return nil
}

func (a *IndexerApp) handleBondProviderEvent(evt types.BondProviderEvent) error {
	provider, err := a.db.FindProvider(evt.Pubkey, evt.Chain)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", evt.Pubkey, evt.Chain)
	}
	if provider == nil {
		// new provider for chain, insert
		if provider, err = a.createProvider(evt); err != nil {
			return errors.Wrapf(err, "error creating provider %s chain %s", evt.Pubkey, evt.Chain)
		}
	} else {
		if evt.BondAbsolute != "" {
			provider.Bond = evt.BondAbsolute
		}
		if _, err = a.db.UpdateProvider(provider); err != nil {
			return errors.Wrapf(err, "error updating provider for bond event %s chain %s", evt.Pubkey, evt.Chain)
		}
	}

	log.Debugf("handled bond provider event for %s chain %s", evt.Pubkey, evt.Chain)
	if _, err = a.db.InsertBondProviderEvent(provider.ID, evt); err != nil {
		return errors.Wrapf(err, "error inserting BondProviderEvent for %s chain %s", evt.Pubkey, evt.Chain)
	}
	return nil
}

func (a *IndexerApp) createProvider(evt types.BondProviderEvent) (*db.ArkeoProvider, error) {
	// new provider for chain, insert
	provider := &db.ArkeoProvider{Pubkey: evt.Pubkey, Chain: evt.Chain, Bond: evt.BondAbsolute}
	entity, err := a.db.InsertProvider(provider)
	if err != nil {
		return nil, errors.Wrapf(err, "error inserting provider %s %s", evt.Pubkey, evt.Chain)
	}
	if entity == nil {
		return nil, fmt.Errorf("nil entity after inserting provider")
	}
	log.Debugf("inserted provider record %d for %s %s", entity.ID, evt.Pubkey, evt.Chain)
	provider.Entity = *entity
	return provider, nil
}

func validateMetadataURI(uri string) bool {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return false
	}
	return true
}
