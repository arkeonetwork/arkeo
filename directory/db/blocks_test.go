package db

import (
	"testing"
	"time"
)

func TestInsertBlock(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}
	b, err := db.InsertBlock(&Block{Height: 1, Hash: "integrationtestblock", BlockTime: time.Now()})
	if err != nil {
		t.Fatalf("error inserting block: %+v", err)
	}
	log.Infof("inserted block b", b)
}

func TestFindLatestBlock(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}
	b, err := db.FindLatestBlock()
	if err != nil {
		t.Fatalf("error: %+v", err)
	}
	log.Infof("found block b %v", b)
}

func TestFindBlockGaps(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}
	b, err := db.FindBlockGaps()
	if err != nil {
		log.Fatalf("error finding gaps: %+v", err)
	}
	log.Infof("gaps: %v", b)
}
