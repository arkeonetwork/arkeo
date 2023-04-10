package db

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common"
	"github.com/stretchr/testify/assert"
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

func TestFindContractsByPubKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	db, err := New(config)
	if err != nil {
		t.Errorf("error getting db: %+v", err)
	}

	provider := "arkeopub1addwnpepqfqqxap0fdehn3jc2vzf2q8hpge2lp65r9gettxev28rqn24mjxtqe9fln3"
	client := "arkeopub1addwnpepq0s3lv5ne868p7jrtcua7awzz9eu76dlddrrpstadrcmtnmyvfk7qsmxt0g"
	// delegate := ""
	contracts, err := db.FindContractsByPubKeys(common.GAIAChainRPCArchiveService.String(), provider, client)
	assert.NoError(t, err)
	assert.NotEmpty(t, contracts)
}
