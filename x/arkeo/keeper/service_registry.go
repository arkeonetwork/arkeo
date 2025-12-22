package keeper

import (
	"fmt"
	"math"
	"strings"

	errorsmod "cosmossdk.io/errors"
	prefix "cosmossdk.io/store/prefix"
	"github.com/arkeonetwork/arkeo/common"
	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// internal helpers to build service registry keys without uppercasing.
func serviceNameKey(name string) []byte {
	return []byte(strings.ToLower(name))
}

func serviceIDKey(id uint64) []byte {
	return []byte(fmt.Sprintf("%s%d", types.ServiceIDKeyPrefix, id))
}

// SetService writes/updates a service record and its secondary id->name index.
func (k KVStore) SetService(ctx cosmos.Context, svc types.Service) error {
	if strings.TrimSpace(svc.Name) == "" {
		return fmt.Errorf("service name cannot be empty")
	}
	if svc.Id == 0 {
		return fmt.Errorf("service id must be non-zero")
	}

	store := ctx.KVStore(k.storeKey)
	nameStore := prefix.NewStore(store, []byte(types.ServiceNameKeyPrefix))

	// enforce unique name
	nameStore.Set(serviceNameKey(svc.Name), k.cdc.MustMarshal(&svc))
	// secondary index: id -> name
	store.Set(serviceIDKey(svc.Id), []byte(strings.ToLower(svc.Name)))

	return nil
}

// GetService looks up a service by name (case-insensitive).
func (k KVStore) GetService(ctx cosmos.Context, name string) (types.Service, bool) {
	store := ctx.KVStore(k.storeKey)
	nameStore := prefix.NewStore(store, []byte(types.ServiceNameKeyPrefix))
	bz := nameStore.Get(serviceNameKey(name))
	if bz == nil {
		return types.Service{}, false
	}
	var svc types.Service
	k.cdc.MustUnmarshal(bz, &svc)
	return svc, true
}

// GetServiceByID looks up a service via its numeric id using the secondary index.
func (k KVStore) GetServiceByID(ctx cosmos.Context, id uint64) (types.Service, bool) {
	store := ctx.KVStore(k.storeKey)
	nameBytes := store.Get(serviceIDKey(id))
	if nameBytes == nil {
		return types.Service{}, false
	}
	return k.GetService(ctx, string(nameBytes))
}

// IterateServices walks all service records, invoking cb for each.
// If cb returns true, iteration stops early.
func (k KVStore) IterateServices(ctx cosmos.Context, cb func(types.Service) bool) {
	store := ctx.KVStore(k.storeKey)
	pStore := prefix.NewStore(store, []byte(types.ServiceNameKeyPrefix))
	it := pStore.Iterator(nil, nil)
	defer it.Close()

	for ; it.Valid(); it.Next() {
		var svc types.Service
		k.cdc.MustUnmarshal(it.Value(), &svc)
		if stop := cb(svc); stop {
			return
		}
	}
}

// RemoveService deletes a service by name and cleans up the secondary id index.
func (k KVStore) RemoveService(ctx cosmos.Context, name string) error {
	store := ctx.KVStore(k.storeKey)
	lowerName := strings.ToLower(name)

	// fetch id for secondary index cleanup
	if svc, ok := k.GetService(ctx, lowerName); ok {
		store.Delete(serviceIDKey(svc.Id))
	}

	nameStore := prefix.NewStore(store, []byte(types.ServiceNameKeyPrefix))
	nameStore.Delete(serviceNameKey(lowerName))
	return nil
}

// ResolveServiceEnum resolves a service name to a common.Service (int32-backed) and returns the registry record.
// It enforces that the id fits in an int32.
func (k KVStore) ResolveServiceEnum(ctx cosmos.Context, name string) (common.Service, types.Service, error) {
	svc, ok := k.GetService(ctx, name)
	if !ok {
		return 0, types.Service{}, errorsmod.Wrapf(sdkerrors.ErrNotFound, "service not found: %s", name)
	}
	if svc.Id == 0 || svc.Id > math.MaxInt32 {
		return 0, types.Service{}, errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid service id for %s: %d", svc.Name, svc.Id)
	}
	return common.Service(int32(svc.Id)), svc, nil
}

// EnsureServiceRegistrySeeded seeds the registry from the legacy static map once.
func (k KVStore) EnsureServiceRegistrySeeded(ctx cosmos.Context) {
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(types.ServiceSeedMarkerKey)) {
		return
	}

	// If the registry already has entries, just mark seeded.
	hasAny := false
	k.IterateServices(ctx, func(types.Service) bool {
		hasAny = true
		return true
	})
	if hasAny {
		store.Set([]byte(types.ServiceSeedMarkerKey), []byte{1})
		return
	}

	for name, id := range common.ServiceLookup {
		desc := common.ServiceDescriptionMap[name]
		_ = k.SetService(ctx, types.Service{
			Id:          uint64(id),
			Name:        name,
			Description: desc,
		})
	}
	store.Set([]byte(types.ServiceSeedMarkerKey), []byte{1})
}
