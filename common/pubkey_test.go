package common

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	. "gopkg.in/check.v1"

	"mercury/common/cosmos"
)

type PubKeyTestSuite struct{}

var _ = Suite(&PubKeyTestSuite{})

// TestPubKey implementation
func (s *PubKeyTestSuite) TestPubKey(c *C) {
	_, pubKey, _ := testdata.KeyTestPubAddr()
	spk, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey)
	c.Assert(err, IsNil)
	pk, err := NewPubKey(spk)
	c.Assert(err, IsNil)
	hexStr := pk.String()
	c.Assert(len(hexStr) > 0, Equals, true)
	pk1, err := NewPubKey(hexStr)
	c.Assert(err, IsNil)
	c.Assert(pk.Equals(pk1), Equals, true)

	result, err := json.Marshal(pk)
	c.Assert(err, IsNil)
	c.Log(result, Equals, fmt.Sprintf(`"%s"`, hexStr))
	var pk2 PubKey
	err = json.Unmarshal(result, &pk2)
	c.Assert(err, IsNil)
	c.Assert(pk2.Equals(pk), Equals, true)
}

func (s *PubKeyTestSuite) TestEquals(c *C) {
	var pk1, pk2, pk3, pk4 PubKey
	_, pubKey1, _ := testdata.KeyTestPubAddr()
	tpk1, err1 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey1)
	c.Assert(err1, IsNil)
	pk1 = PubKey(tpk1)

	_, pubKey2, _ := testdata.KeyTestPubAddr()
	tpk2, err2 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey2)
	c.Assert(err2, IsNil)
	pk2 = PubKey(tpk2)

	_, pubKey3, _ := testdata.KeyTestPubAddr()
	tpk3, err3 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey3)
	c.Assert(err3, IsNil)
	pk3 = PubKey(tpk3)

	_, pubKey4, _ := testdata.KeyTestPubAddr()
	tpk4, err4 := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pubKey4)
	c.Assert(err4, IsNil)
	pk4 = PubKey(tpk4)

	c.Assert(PubKeys{
		pk1, pk2,
	}.Equals(nil), Equals, false)

	c.Assert(PubKeys{
		pk1, pk2, pk3,
	}.Equals(PubKeys{
		pk1, pk2,
	}), Equals, false)

	c.Assert(PubKeys{
		pk1, pk2, pk3, pk4,
	}.Equals(PubKeys{
		pk4, pk3, pk2, pk1,
	}), Equals, true)

	c.Assert(PubKeys{ // nolint
		pk1, pk2, pk3, pk4,
	}.Equals(PubKeys{
		pk1, pk2, pk3, pk4,
	}), Equals, true)
}
