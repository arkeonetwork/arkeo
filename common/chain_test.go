package common

import (
	. "gopkg.in/check.v1"
)

type ChainSuite struct{}

var _ = Suite(&ChainSuite{})

func (s ChainSuite) TestChain(c *C) {
	bnbChain, err := NewChain("mcy")
	c.Assert(err, IsNil)
	c.Check(bnbChain.Equals(BaseChain), Equals, true)
	c.Check(bnbChain.IsEmpty(), Equals, false)
	c.Check(bnbChain.String(), Equals, "MCY")

	_, err = NewChain("B") // too short
	c.Assert(err, NotNil)

	chains := Chains{"BNB", "BTC"}
	c.Check(chains.Has("BTC"), Equals, true)
	c.Check(chains.Has("ETH"), Equals, false)
}
