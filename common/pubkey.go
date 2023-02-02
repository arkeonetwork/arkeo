package common

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/arkeonetwork/arkeo/common/cosmos"

	"github.com/btcsuite/btcutil/bech32"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types"
)

type (
	PubKey  string
	PubKeys []PubKey
)

var EmptyPubKey PubKey

// NewPubKey create a new instance of PubKey
// key is bech32 encoded string
func NewPubKey(key string) (PubKey, error) {
	if len(key) == 0 {
		return EmptyPubKey, nil
	}
	_, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, key)
	if err != nil {
		return EmptyPubKey, fmt.Errorf("%s is not bech32 encoded pub key,err : %w", key, err)
	}
	return PubKey(key), nil
}

func NewPubKeyFromCrypto(pk cryptotypes.PubKey) (PubKey, error) {
	/*
		tmp, err := codec.ToTmPubKeyInterface(pk)
		if err != nil {
			return EmptyPubKey, fmt.Errorf("fail to create PubKey from crypto.PubKey,err:%w", err)
		}
	*/
	s, err := cosmos.Bech32ifyPubKey(cosmos.Bech32PubKeyTypeAccPub, pk)
	if err != nil {
		return EmptyPubKey, fmt.Errorf("fail to create PubKey from crypto.PubKey,err:%w", err)
	}
	return PubKey(s), nil
}

// Equals check whether two are the same
func (pubKey PubKey) Equals(pubKey1 PubKey) bool {
	return pubKey == pubKey1
}

// IsEmpty to check whether it is empty
func (pubKey PubKey) IsEmpty() bool {
	return len(pubKey) == 0
}

// String stringer implementation
func (pubKey PubKey) String() string {
	return string(pubKey)
}

func (pubKey PubKey) GetMyAddress() (cosmos.AccAddress, error) {
	pk, err := cosmos.GetPubKeyFromBech32(cosmos.Bech32PubKeyTypeAccPub, string(pubKey))
	if err != nil {
		return cosmos.AccAddress{}, err
	}
	prefix := types.GetConfig().GetBech32AccountAddrPrefix()
	addr, err := ConvertAndEncode(prefix, pk.Address().Bytes())
	if err != nil {
		return cosmos.AccAddress{}, fmt.Errorf("fail to bech32 encode the address, err: %w", err)
	}
	return cosmos.AccAddressFromBech32(addr)
}

// MarshalJSON to Marshals to JSON using Bech32
func (pubKey PubKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(pubKey.String())
}

// UnmarshalJSON to Unmarshal from JSON assuming Bech32 encoding
func (pubKey *PubKey) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	pk, err := NewPubKey(s)
	if err != nil {
		return err
	}
	*pubKey = pk
	return nil
}

func (pks PubKeys) Valid() error {
	for _, pk := range pks {
		if _, err := NewPubKey(pk.String()); err != nil {
			return err
		}
	}
	return nil
}

func (pks PubKeys) Contains(pk PubKey) bool {
	for _, p := range pks {
		if p.Equals(pk) {
			return true
		}
	}
	return false
}

// Equals check whether two pub keys are identical
func (pks PubKeys) Equals(newPks PubKeys) bool {
	if len(pks) != len(newPks) {
		return false
	}

	source := make(PubKeys, len(pks))
	dest := make(PubKeys, len(newPks))
	copy(source, pks)
	copy(dest, newPks)

	// sort both lists
	sort.Slice(source[:], func(i, j int) bool {
		return source[i].String() < source[j].String()
	})
	sort.Slice(dest[:], func(i, j int) bool {
		return dest[i].String() < dest[j].String()
	})
	for i := range source {
		if !source[i].Equals(dest[i]) {
			return false
		}
	}
	return true
}

// String implement stringer interface
func (pks PubKeys) String() string {
	strs := make([]string, len(pks))
	for i := range pks {
		strs[i] = pks[i].String()
	}
	return strings.Join(strs, ", ")
}

func (pks PubKeys) Strings() []string {
	allStrings := make([]string, len(pks))
	for i, pk := range pks {
		allStrings[i] = pk.String()
	}
	return allStrings
}

// ConvertAndEncode converts from a base64 encoded byte string to hex or base32 encoded byte string and then to bech32
func ConvertAndEncode(hrp string, data []byte) (string, error) {
	converted, err := bech32.ConvertBits(data, 8, 5, true)
	if err != nil {
		return "", fmt.Errorf("encoding bech32 failed,%w", err)
	}
	return bech32.Encode(hrp, converted)
}
