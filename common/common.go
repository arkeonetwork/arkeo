package common

import (
	"arkeo/common/cosmos"
	"fmt"
)

// GetSafeShare does the same as GetUncappedShare , but GetSafeShare will guarantee the result will not more than total
func GetSafeShare(part, total, allocation cosmos.Int) cosmos.Int {
	if part.GTE(total) {
		part = total
	}
	return GetUncappedShare(part, total, allocation)
}

// GetUncappedShare this method will panic if any of the input parameter can't be convert to cosmos.Dec
// which shouldn't happen
func GetUncappedShare(part, total, allocation cosmos.Int) (share cosmos.Int) {
	if part.IsZero() || total.IsZero() {
		return cosmos.ZeroInt()
	}
	defer func() {
		if err := recover(); err != nil {
			share = cosmos.ZeroInt()
		}
	}()
	// use string to convert cosmos.Int to cosmos.Dec is the only way I can find out without being constrain to uint64
	// cosmos.Int can hold values way larger than uint64 , because it is using big.Int internally
	aD, err := cosmos.NewDecFromStr(allocation.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", allocation.String(), err))
	}

	pD, err := cosmos.NewDecFromStr(part.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert %s to cosmos.Dec: %w", part.String(), err))
	}
	tD, err := cosmos.NewDecFromStr(total.String())
	if err != nil {
		panic(fmt.Errorf("fail to convert%s to cosmos.Dec: %w", total.String(), err))
	}
	// A / (Total / part) == A * (part/Total) but safer when part < Totals
	result := aD.Quo(tD.Quo(pD))
	share = cosmos.NewIntFromBigInt(result.RoundInt().BigInt())
	return
}

func Tokens(i int64) int64 {
	return i * 1e8
}
