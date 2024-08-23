package common

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

func TestGetUncappedShare(t *testing.T) {
	part := cosmos.NewInt(149506590)
	total := cosmos.NewInt(50165561086)
	alloc := cosmos.NewInt(50000000)
	share := GetUncappedShare(part, total, alloc)
	require.True(t, share.Equal(cosmos.NewInt(149013)))
}

func TestGetSafeShare(t *testing.T) {
	part := cosmos.NewInt(14950659000000000)
	total := cosmos.NewInt(50165561086)
	alloc := cosmos.NewInt(50000000)
	share := GetSafeShare(part, total, alloc)
	require.True(t, share.Equal(cosmos.NewInt(50000000)))
}
