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
	delegatePubkey := "arkeopub1addwnpepqglj743j5pchx57g4rwxvlfrgy2mztwq837hu90mrdxmqv09hagrunus4ja"
	providerID := int64(2)
	contract, err := db.FindContract(providerID, delegatePubkey, 0)
	if err != nil {
		t.Errorf("error finding contract: %+v", err)
		t.FailNow()
	}
	log.Infof("found contract %d", contract.ID)

	delegatePubkey = "nosuchthing"
	contract, err = db.FindContract(providerID, delegatePubkey, 0)
	if err != nil {
		t.Errorf("error finding contract: %+v", err)
		t.FailNow()
	}
	if contract != nil {
		t.Errorf("expected nil but got %v", contract)
	}
}
