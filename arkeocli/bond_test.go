package arkeocli

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common/cosmos"
	"github.com/stretchr/testify/require"
)

func TestParseBondAmount(t *testing.T) {
	amount, err := parseBondAmount("100uarkeo")
	require.NoError(t, err)
	require.Equal(t, cosmos.NewInt(100), amount)

	amount, err = parseBondAmount("-100uarkeo")
	require.NoError(t, err)
	require.Equal(t, cosmos.NewInt(-100), amount)

	_, err = parseBondAmount("100arkeo")
	require.ErrorContains(t, err, "bad bond denom")

	_, err = parseBondAmount("-100arkeo")
	require.ErrorContains(t, err, "bad bond denom")

	_, err = parseBondAmount("100")
	require.Error(t, err)
	require.ErrorContains(t, err, "bad bond denom")

	_, err = parseBondAmount("-100")
	require.ErrorContains(t, err, "bad bond denom")

	_, err = parseBondAmount("")
	require.ErrorContains(t, err, "bad bond denom")

	_, err = parseBondAmount("uarkeo")
	require.ErrorContains(t, err, "bad bond amount")

	_, err = parseBondAmount("uarkeo")
	require.ErrorContains(t, err, "bad bond amount")

	_, err = parseBondAmount("Paul Revere")
	require.ErrorContains(t, err, "bad bond denom")

	_, err = parseBondAmount("100uarkeo 100uarkeo")
	require.ErrorContains(t, err, "bad bond denom")
}
