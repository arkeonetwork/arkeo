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
	"github.com/arkeonetwork/arkeo/x/arkeo/configs"
	arkeoTypes "github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/pkg/errors"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func (a *IndexerApp) handleModProviderEvent(evt ctypes.ResultEvent) error {
	typedEvent, err := arkeoUtils.ParseTypedEvent(evt, "arkeo.arkeo.EventModProvider")
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
	modProviderEvent := typedEvent.(*arkeoTypes.EventModProvider)
	provider, err := a.db.FindProvider(modProviderEvent.Provider.String(), modProviderEvent.Service)
	if err != nil {
		return errors.Wrapf(err, "error finding provider %s for chain %s", modProviderEvent.Provider.String(), modProviderEvent.Service)
	}
	if provider == nil {
		return fmt.Errorf("cannot mod provider, DNE %s %s", modProviderEvent.Provider, modProviderEvent.Service)
	}

	log := log.WithField("provider", strconv.FormatInt(provider.ID, 10))

	isMetaDataUpdated := provider.MetadataNonce == 0 || provider.MetadataNonce < modProviderEvent.MetadataNonce
	provider.MetadataURI = modProviderEvent.MetadataURI
	provider.MetadataNonce = modProviderEvent.MetadataNonce
	provider.Status = types.ProviderStatus(modProviderEvent.Status.String())
	provider.MinContractDuration = modProviderEvent.MinContractDuration
	provider.MaxContractDuration = modProviderEvent.MaxContractDuration

	if len(modProviderEvent.SubscriptionRate) == 1 && modProviderEvent.SubscriptionRate[0].Denom == configs.Denom {
		provider.SubscriptionRate = modProviderEvent.SubscriptionRate[0].Amount.Int64()
	} else {
		log.Warnf("only a single subscription rate with denom %s is supported", configs.Denom)
	}

	if len(modProviderEvent.PayAsYouGoRate) == 1 && modProviderEvent.PayAsYouGoRate[0].Denom == configs.Denom {
		provider.PayAsYouGoRate = modProviderEvent.PayAsYouGoRate[0].Amount.Int64()
	} else {
		log.Warnf("only a single pay as you go rate with denom %s is supported", configs.Denom)
	}

	if _, err = a.db.UpdateProvider(provider); err != nil {
		return errors.Wrapf(err, "error updating provider for mod event %s chain %s", provider.Pubkey, provider.Chain)
	}
	log.Infof("updated provider %s chain %s", provider.Pubkey, provider.Chain)

	mpe := types.ModProviderEvent{
		Pubkey:              modProviderEvent.Provider.String(),
		Chain:               modProviderEvent.Service,
		Height:              height,
		TxID:                txid,
		MetadataURI:         modProviderEvent.MetadataURI,
		MetadataNonce:       modProviderEvent.MetadataNonce,
		Status:              types.ProviderStatus(modProviderEvent.Status.String()),
		MinContractDuration: modProviderEvent.MinContractDuration,
		MaxContractDuration: modProviderEvent.MaxContractDuration,
		SubscriptionRate:    provider.SubscriptionRate,
		PayAsYouGoRate:      provider.PayAsYouGoRate,
	}
	if _, err = a.db.InsertModProviderEvent(provider.ID, mpe); err != nil {
		return errors.Wrapf(err, "error inserting ModProviderEvent for %s chain %s", modProviderEvent.Provider, modProviderEvent.Service)
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
