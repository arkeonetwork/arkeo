package db

import (
	"testing"
)

func TestGetArkeoNetworkStats(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}

	stats, err := db.GetArkeoNetworkStats()
	if err != nil {
		t.Error("error getting stats", err)
		t.FailNow()
	}

	if stats.ContractsTotal == 0 {
		t.Error("error getting stats with data", err)
		t.FailNow()
	}
}
