package indexer

import (
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	arkeoUtils "github.com/arkeonetwork/arkeo/common/utils"
	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/directory/utils"
	arkeoTypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
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

func (a *IndexerApp) handleBondProviderEvent(evt ctypes.ResultEvent) error {
	typedEvent, err := arkeoUtils.ParseTypedEvent(evt, "arkeo.arkeo.EventBondProvider")
	if err != nil {
		log.Errorf("failed to parse typed event", "error", err)
		return errors.Wrapf(err, "failed to parse typed event")
	}

	txData, ok := evt.Data.(tmtypes.EventDataTx)
	if !ok {
		return fmt.Errorf("failed to cast %T to EventDataTx", evt.Data)
	}

	height := txData.Height
	txid := strings.ToUpper(hex.EncodeToString(tmtypes.Tx(txData.Tx).Hash()))

	bondProviderEvent, ok := typedEvent.(*arkeoTypes.EventBondProvider)
	if !ok {
		return fmt.Errorf("failed to cast %T to EventBondProvider", typedEvent)
	}

	provider, err := a.db.FindProvider(bondProviderEvent.Provider.String(), bondProviderEvent.Service)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", bondProviderEvent.Provider, bondProviderEvent.Service)
	}
	if provider == nil {
		// new provider for chain, insert
		if provider, err = a.createProvider(bondProviderEvent); err != nil {
			return errors.Wrapf(err, "error creating provider %s chain %s", bondProviderEvent.Provider, bondProviderEvent.Service)
		}
	} else {
		provider.Bond = bondProviderEvent.BondAbs.String()
		if _, err = a.db.UpdateProvider(provider); err != nil {
			return errors.Wrapf(err, "error updating provider for bond event %s chain %s", bondProviderEvent.Provider, bondProviderEvent.Service)
		}
	}

	log.Debugf("handled bond provider event for %s chain %s", bondProviderEvent.Provider.String(), bondProviderEvent.Service)
	bpe := types.BondProviderEvent{
		Pubkey:       bondProviderEvent.Provider.String(),
		Chain:        bondProviderEvent.Service,
		Height:       height,
		TxID:         txid,
		BondRelative: bondProviderEvent.BondRel.String(),
		BondAbsolute: bondProviderEvent.BondAbs.String(),
	}

	if _, err = a.db.InsertBondProviderEvent(provider.ID, bpe); err != nil {
		return errors.Wrapf(err, "error inserting BondProviderEvent for %s chain %s", bondProviderEvent.Provider.String(), bondProviderEvent.Service)
	}
	return nil
}

func (a *IndexerApp) createProvider(evt *arkeoTypes.EventBondProvider) (*db.ArkeoProvider, error) {
	// new provider for chain, insert
	provider := &db.ArkeoProvider{Pubkey: evt.Provider.String(), Chain: evt.Service, Bond: evt.BondAbs.String()}
	entity, err := a.db.InsertProvider(provider)
	if err != nil {
		return nil, errors.Wrapf(err, "error inserting provider %s %s", evt.Provider.String(), evt.Service)
	}
	if entity == nil {
		return nil, fmt.Errorf("nil entity after inserting provider")
	}
	log.Debugf("inserted provider record %d for %s %s", entity.ID, evt.Provider.String(), evt.Service)
	provider.Entity = *entity
	return provider, nil
}

func validateMetadataURI(uri string) bool {
	if _, err := url.ParseRequestURI(uri); err != nil {
		return false
	}
	return true
}
