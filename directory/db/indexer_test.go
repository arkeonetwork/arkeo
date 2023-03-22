package db

import (
	"testing"
)

func TestInsertIndexerStatus(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}

	entity, err := db.UpsertIndexerStatus(&IndexerStatus{
		ID:     0,
		Height: 55,
	})

	if err != nil {
		t.Errorf("error inserting indexer: %+v", err)
		t.FailNow()
	}
	log.Infof("inserted indexer %d", entity.ID)
}

func TestUpdateIndexerStatus(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}

	entity, err := db.UpdateIndexerStatus(&IndexerStatus{
		ID:     0,
		Height: 65,
	})

	if err != nil {
		t.Errorf("error updated indexer: %+v", err)
		t.FailNow()
	}
	log.Infof("updated indexer %d", entity.ID)
}

func TestFindIndexerStatus(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}

	var id int64
	id = 0
	indexerStatus, err := db.FindIndexerStatus(id)
	if err != nil {
		t.Errorf("error finding indexer status: %+v", err)
		t.FailNow()
	}
	log.Infof("found indexer status %d", indexerStatus.ID)

	id = 555556
	indexerStatus, err = db.FindIndexerStatus(id)
	if err != nil {
		t.Errorf("error finding provider: %+v", err)
		t.FailNow()
	}
	if indexerStatus != nil {
		t.Errorf("expected nil but got %v", indexerStatus)
	}
}
