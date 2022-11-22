package common

import (
	. "gopkg.in/check.v1"
)

type ChainSuite struct{}

var _ = Suite(&ChainSuite{})

func (s ChainSuite) TestChain(c *C) {
	chn, err := NewChain("btc-mainnet-fullnode")
	c.Assert(err, IsNil)
	c.Check(chn.Equals(BTCChain), Equals, true)
	c.Check(chn.IsEmpty(), Equals, false)
	c.Check(chn.String(), Equals, "btc-mainnet-fullnode")

	_, err = NewChain("B") // invalid
	c.Assert(err, NotNil)
}
