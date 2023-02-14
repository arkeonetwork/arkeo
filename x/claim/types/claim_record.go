package types

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
