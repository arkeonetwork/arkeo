package types

import (
	"fmt"
	"strings"
)

func (claimRecord *ClaimRecord) IsEmpty() bool {
	if *claimRecord == (ClaimRecord{}) {
		return true
	}

	if claimRecord.Address == "" {
		return true
	}

	if !claimRecord.AmountClaim.IsNil() && claimRecord.AmountClaim.IsValid() && !claimRecord.AmountClaim.IsZero() {
		return false
	}

	if !claimRecord.AmountVote.IsNil() && claimRecord.AmountVote.IsValid() && !claimRecord.AmountVote.IsZero() {
		return false
	}

	if !claimRecord.AmountDelegate.IsNil() && claimRecord.AmountDelegate.IsValid() && !claimRecord.AmountDelegate.IsZero() {
		return false
	}

	return true
}

// ChainFromString convert chain string to Chain Enum type
func ChainFromString(chain string) (Chain, error) {
	for id, item := range Chain_name {
		if strings.EqualFold(item, chain) {
			return Chain(id), nil
		}
	}
	return Chain(-1), fmt.Errorf("invalid chain")
}
