package db

import (
	"testing"

	"github.com/pashagolub/pgxmock/v2"
	"github.com/stretchr/testify/assert"
)

func TestGetArkeoNetworkStats(t *testing.T) {
	m, db := getMockDirectoryDBForTest(t)
	defer m.Close()
	m.ExpectQuery("select.*from network_stats_v.*").
		WillReturnRows(
			pgxmock.NewRows([]string{
				"open_contracts", "total_contracts", "median_open_contract_length",
				"median_open_contract_rate", "total_online_providers", "total_queries", "total_paid",
			}).
				AddRow(int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7)))
	state, err := db.GetArkeoNetworkStats()
	assert.Nil(t, err)
	assert.NotNil(t, state)
	assert.Equal(t, int64(1), state.ContractsOpen)
	assert.Equal(t, int64(2), state.ContractsTotal)
	assert.Equal(t, int64(3), state.ContractsMedianDuration)
	assert.Equal(t, int64(4), state.ContractsMedianRate)
	assert.Equal(t, int64(5), state.ProviderCount)
	assert.Equal(t, int64(6), state.QueryCount)
	assert.Equal(t, int64(7), state.TotalIncome)
	assert.Nil(t, m.ExpectationsWereMet())
}
