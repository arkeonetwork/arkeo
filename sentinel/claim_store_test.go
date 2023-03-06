package sentinel

import (
	"io/ioutil"
	"os"

	"github.com/arkeonetwork/arkeo/x/arkeo/types"

	. "gopkg.in/check.v1"
)

type ClaimStoreSuite struct {
	dir string
}

var _ = Suite(&ClaimStoreSuite{})

func (s *ClaimStoreSuite) SetUpTest(c *C) {
	var err error
	s.dir, err = ioutil.TempDir("/tmp", "claim-store")
	c.Assert(err, IsNil)
}

func (s *ClaimStoreSuite) TestStore(c *C) {
	store, err := NewClaimStore(s.dir)
	c.Assert(err, IsNil)

	pk2 := types.GetRandomPubKey()
	contractId := uint64(57)
	c.Assert(err, IsNil)
	claim := NewClaim(contractId, pk2, 30, 10, "signature")

	c.Assert(store.Set(claim), IsNil)
	c.Assert(store.Has(claim.Key()), Equals, true)
	claim, err = store.Get(claim.Key())
	c.Assert(err, IsNil)
	c.Check(claim.Height, Equals, int64(10))

	claims := store.List()
	c.Assert(claims, HasLen, 1)

	c.Assert(store.Remove(claim.Key()), IsNil)
	c.Assert(store.Has(claim.Key()), Equals, false)
}

func (s *ClaimStoreSuite) TearDownSuite(c *C) {
	defer os.RemoveAll(s.dir)
}
