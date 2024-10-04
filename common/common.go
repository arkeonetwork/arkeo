package common

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"

	"github.com/arkeonetwork/arkeo/common/cosmos"
)

const (
	// SuperMajorityFactor - super majority 2/3
	SuperMajorityFactor = 3
	// SimpleMajorityFactor - simple majority 1/2
	SimpleMajorityFactor = 2
)

// GetSafeShare does the same as GetUncappedShare , but GetSafeShare will guarantee the result will not more than total
func GetSafeShare(part, total, allocation cosmos.Dec) cosmos.Dec {
	if part.GTE(total) {
		part = total
	}
	return GetUncappedShare(part, total, allocation)
}

// GetUncappedShare this method will panic if any of the input parameter can't be convert to cosmos.Dec
// which shouldn't happen
func GetUncappedShare(part, total, allocation cosmos.Dec) (share cosmos.Dec) {
	if part.IsZero() || total.IsZero() {
		return cosmos.ZeroDec()
	}
	defer func() {
		if err := recover(); err != nil {
			share = cosmos.ZeroDec()
		}
	}()
	// A / (Total / part) == A * (part/Total) but safer when part < Totals
	share = allocation.Quo(total.Quo(part))
	return
}

func Tokens(i int64) int64 {
	return i * 1e8
}

func MustParseURL(uri string) *url.URL {
	uri = strings.TrimSpace(uri)
	u, err := url.Parse(uri)
	if err != nil {
		panic(fmt.Errorf("unable to parse uri %s: %w", uri, err))
	}
	return u
}

// GetCurrentVersion - intended for unit tests, fetches the current version of
// arkeo via `version` file
// #nosec G304 this is a method only used for test purpose
func GetCurrentVersion() int64 {
	_, filename, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filename), "./../")
	dat, err := os.ReadFile(path.Join(dir, "chain.version"))
	if err != nil {
		panic(err)
	}
	v, err := strconv.ParseInt(strings.TrimSpace(string(dat)), 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}
