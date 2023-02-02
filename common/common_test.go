package common

import (
	"testing"

	"github.com/arkeonetwork/arkeo/common/cosmos"

	. "gopkg.in/check.v1"
)

func TestPackage(t *testing.T) { TestingT(t) }

type CommonSuite struct{}

var _ = Suite(&CommonSuite{})

func (s CommonSuite) TestGetUncappedShare(c *C) {
	part := cosmos.NewInt(149506590)
	total := cosmos.NewInt(50165561086)
	alloc := cosmos.NewInt(50000000)
	share := GetUncappedShare(part, total, alloc)
	c.Assert(share.Equal(cosmos.NewInt(149013)), Equals, true)
}

func (s CommonSuite) TestGetSafeShare(c *C) {
	part := cosmos.NewInt(14950659000000000)
	total := cosmos.NewInt(50165561086)
	alloc := cosmos.NewInt(50000000)
	share := GetSafeShare(part, total, alloc)
	c.Assert(share.Equal(cosmos.NewInt(50000000)), Equals, true)
}
