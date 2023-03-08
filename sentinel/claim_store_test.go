package sentinel

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ClaimStoreSuite struct {
	suite.Suite
	dir string
}

func (s *ClaimStoreSuite) SetUpTest() {
	var err error
	s.dir, err = ioutil.TempDir("/tmp", "claim-store")
	require.NoError(s.T(), err)
}

func (s *ClaimStoreSuite) TestStore() {
	store, err := NewClaimStore(s.dir)
	require.NoError(s.T(), err)

	pk2 := types.GetRandomPubKey()
	contractId := uint64(57)
	require.NoError(s.T(), err)
	claim := NewClaim(contractId, pk2, 30, "signature")

	require.Nil(s.T(), store.Set(claim))

	require.True(s.T(), store.Has(claim.Key()))
	claim, err = store.Get(claim.Key())
	require.NoError(s.T(), err)

	claims := store.List()
	require.Len(s.T(), claims, 1)

	require.NoError(s.T(), store.Remove(claim.Key()))
	require.False(s.T(), store.Has(claim.Key()))
}

func (s *ClaimStoreSuite) TearDownSuite(t *testing.T) {
	defer os.RemoveAll(s.dir)
}

func TestClaimStoreSuite(t *testing.T) {
	suite.Run(t, new(ClaimStoreSuite))
}
