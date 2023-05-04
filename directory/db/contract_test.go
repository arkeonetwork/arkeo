package db

import (
	"testing"
)

func TestFindContract(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}
	contractId := uint64(2)
	contract, err := db.FindContract(contractId)
	if err != nil {
		t.Errorf("error finding contract: %+v", err)
		t.FailNow()
	}
	log.Infof("found contract %d", contract.ID)

	contract, err = db.FindContract(uint64(5))
	if err != nil {
		t.Errorf("error finding contract: %+v", err)
		t.FailNow()
	}
	if contract != nil {
		t.Errorf("expected nil but got %v", contract)
	}
}
