package sentinel

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
)

func TestNewProviderConfigurationStore_InMemory(t *testing.T) {
	store, err := NewProviderConfigurationStore("")
	require.NoError(t, err, "expected no error when creating in-memory store")
	require.NotNil(t, store, "expected store to be created")

	err = store.db.Close()
	require.NoError(t, err, "expected no error when closing in-memory store")
}

func TestNewProviderConfigurationStore_FileBased(t *testing.T) {
	tempDir := t.TempDir()
	store, err := NewProviderConfigurationStore(tempDir)
	require.NoError(t, err, "expected no error when creating file-based store")
	require.NotNil(t, store, "expected store to be created")

	err = store.db.Close()
	require.NoError(t, err, "expected no error when closing file-based store")
}

func TestProviderConfigurationStore_SetAndGet(t *testing.T) {
	store, err := NewProviderConfigurationStore("")
	require.NoError(t, err, "expected no error when creating in-memory store")

	pubKey, _ := common.NewPubKey("tarkeopub1addwnpepqfzke9590mrh4m430zapyl3eh0na4ffzrssz89d4qq89ffuy4xn2yqgcm5v")
	service := common.Service(common.ServiceLookup["mock"])
	config := ProviderConfiguration{
		PubKey:              pubKey,
		Service:             service,
		Bond:                cosmos.NewInt(1000),
		BondRelative:        cosmos.NewInt(100),
		MetadataUri:         "http://test-metadata.com",
		MetadataNonce:       1,
		Status:              types.ProviderStatus_ONLINE,
		MinContractDuration: 10,
		MaxContractDuration: 100,
		SubscriptionRate:    cosmos.Coins{cosmos.NewInt64Coin("arkeo", 50)},
		PayAsYouGoRate:      cosmos.Coins{cosmos.NewInt64Coin("arkeo", 5)},
		SettlementDuration:  200,
	}

	err = store.Set(config)
	require.NoError(t, err, "expected no error when setting provider config")

	retrievedConfig, err := store.Get(pubKey, service.String())
	require.NoError(t, err, "expected no error when getting provider config")
	require.Equal(t, config, retrievedConfig, "expected the stored and retrieved configs to be the same")

	err = store.db.Close()
	require.NoError(t, err, "expected no error when closing in-memory store")
}

func TestProviderConfigurationStore_Remove(t *testing.T) {
	store, err := NewProviderConfigurationStore("")
	require.NoError(t, err, "expected no error when creating in-memory store")

	pubKey, _ := common.NewPubKey("tarkeopub1addwnpepqfzke9590mrh4m430zapyl3eh0na4ffzrssz89d4qq89ffuy4xn2yqgcm5v")
	service := common.Service(common.ServiceLookup["mock"])
	config := ProviderConfiguration{
		PubKey:              pubKey,
		Service:             service,
		Bond:                cosmos.NewInt(1000),
		BondRelative:        cosmos.NewInt(100),
		MetadataUri:         "http://test-metadata.com",
		MetadataNonce:       1,
		Status:              types.ProviderStatus_ONLINE,
		MinContractDuration: 10,
		MaxContractDuration: 100,
		SubscriptionRate:    cosmos.Coins{cosmos.NewInt64Coin("arkeo", 50)},
		PayAsYouGoRate:      cosmos.Coins{cosmos.NewInt64Coin("arkeo", 5)},
		SettlementDuration:  200,
	}

	err = store.Set(config)
	require.NoError(t, err, "expected no error when setting provider config")

	err = store.Remove(pubKey, service.String())
	require.NoError(t, err, "expected no error when removing provider config")

	_, err = store.Get(pubKey, service.String())
	require.Error(t, err, "expected error when getting non-existent provider config")

	err = store.db.Close()
	require.NoError(t, err, "expected no error when closing in-memory store")
}
