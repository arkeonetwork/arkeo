package db

import (
	"testing"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/sentinel"
	"github.com/arkeonetwork/arkeo/sentinel/conf"
	"github.com/google/uuid"
)

func TestInsertProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}

	entity, err := db.InsertProvider(&ArkeoProvider{
		Pubkey:  uuid.NewString(),
		Service: "btc-mainnet-fullnode",
		Bond:    "1234567890",
	})
	if err != nil {
		t.Errorf("error inserting provider: %+v", err)
		t.FailNow()
	}
	log.Infof("inserted provider %d", entity.ID)
}

func TestFindProvider(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}
	pubkey := "arkeopub1addwnpepqg5fsc756nx3wlrp7f4328slhgfulhu53epxnyy4q6ln3htrhxxsczgwfyf"
	service := "btc-mainnet"
	provider, err := db.FindProvider(pubkey, service)
	if err != nil {
		t.Errorf("error finding provider: %+v", err)
		t.FailNow()
	}
	log.Infof("found provider %d", provider.ID)

	pubkey = "nosuchthing"
	provider, err = db.FindProvider(pubkey, service)
	if err != nil {
		t.Errorf("error finding provider: %+v", err)
		t.FailNow()
	}
	if provider != nil {
		t.Errorf("expected nil but got %v", provider)
	}
}

func TestUpsertProviderMetadata(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}
	if _, err = db.UpsertProviderMetadata(1, 1, sentinel.Metadata{Version: "0.0.6t", Configuration: conf.Configuration{Moniker: "UnitTestOper", Location: "50.1535,-19.165"}}); err != nil {
		t.Errorf("error upserting: %+v", err)
	}
}

func TestSearchProviders(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}

	searchParams := types.ProviderSearchParams{IsMaxDistanceSet: true, Coordinates: types.Coordinates{Latitude: 50.01, Longitude: -35.68}, MaxDistance: 0}
	results, err := db.SearchProviders(searchParams)
	if err != nil {
		t.Errorf("error finding provider with geolocation: %+v", err)
		t.FailNow()
	}

	if results == nil {
		t.FailNow()
	}

	if len(results) != 0 {
		t.FailNow()
	}

	searchParams.MaxDistance = 1000 // miles
	results, err = db.SearchProviders(searchParams)

	if err != nil {
		t.Errorf("error finding provider with geolocation: %+v", err)
		t.FailNow()
	}

	if results == nil {
		t.FailNow()
	}

	if len(results) < 1 {
		t.FailNow()
	}
}
