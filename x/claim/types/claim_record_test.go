package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestIsEmpty tests IsEmpty method
func TestIsEmpty(t *testing.T) {
	// returns true if empty or nil
	require.True(t, (&ClaimRecord{}).IsEmpty())
	require.True(t, (&ClaimRecord{Address: ""}).IsEmpty())
	require.True(t, (&ClaimRecord{Address: "foo"}).IsEmpty())
	require.True(t, (&ClaimRecord{Address: "foo", AmountClaim: types.Coin{}}).IsEmpty())

	// returns false if not empty
	require.False(t, (&ClaimRecord{Address: "foo", AmountClaim: types.NewInt64Coin("foo", 1)}).IsEmpty())
	require.False(t, (&ClaimRecord{Address: "foo", AmountVote: types.NewInt64Coin("foo", 1)}).IsEmpty())
	require.False(t, (&ClaimRecord{Address: "foo", AmountDelegate: types.NewInt64Coin("foo", 1)}).IsEmpty())
}
