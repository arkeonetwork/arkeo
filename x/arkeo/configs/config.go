package configs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	GitCommit             = "null" // sha1 revision used to build the program
	BuildTime             = "null" // when the executable was built
	Version               = "1"    // software version
	Denom                 = "uarkeo"
	MaxBasisPoints  int64 = 10_000
	int64Overrides        = map[ConfigName]int64{}
	boolOverrides         = map[ConfigName]bool{}
	stringOverrides       = map[ConfigName]string{}
)

// var BlockTime = 5 * time.Second

// ConfigVals implement ConfigValues interface
type ConfigVals struct {
	int64values  map[ConfigName]int64
	boolValues   map[ConfigName]bool
	stringValues map[ConfigName]string
}

// GetInt64Value get value in int64 type, if it doesn't exist then it will return the default value of int64, which is 0
func (cv *ConfigVals) GetInt64Value(name ConfigName) int64 {
	// check overrides first
	v, ok := int64Overrides[name]
	if ok {
		return v
	}

	v, ok = cv.int64values[name]
	if !ok {
		return 0
	}
	return v
}

// GetBoolValue retrieve a bool constant value from the map
func (cv *ConfigVals) GetBoolValue(name ConfigName) bool {
	v, ok := boolOverrides[name]
	if ok {
		return v
	}
	v, ok = cv.boolValues[name]
	if !ok {
		return false
	}
	return v
}

// GetStringValue retrieve a string const value from the map
func (cv *ConfigVals) GetStringValue(name ConfigName) string {
	v, ok := stringOverrides[name]
	if ok {
		return v
	}
	v, ok = cv.stringValues[name]
	if ok {
		return v
	}
	return ""
}

func (cv *ConfigVals) String() string {
	// get all the keys
	int64Keys := make([]ConfigName, 0, len(cv.int64values))
	for k := range cv.int64values {
		int64Keys = append(int64Keys, k)
	}
	sort.Slice(int64Keys, func(i, j int) bool {
		return int64Keys[i].String() < int64Keys[j].String()
	})

	boolKeys := make([]ConfigName, 0, len(cv.boolValues))
	for k := range cv.boolValues {
		boolKeys = append(boolKeys, k)
	}
	sort.Slice(boolKeys, func(i, j int) bool {
		return boolKeys[i].String() < boolKeys[j].String()
	})

	sb := strings.Builder{}
	for _, k := range int64Keys {
		if overrideValue, ok := int64Overrides[k]; ok {
			sb.WriteString(fmt.Sprintf("%s:%d\n", k, overrideValue))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s:%d\n", k, cv.int64values[k]))
	}

	for _, k := range boolKeys {
		if overrideValue, ok := boolOverrides[k]; ok {
			sb.WriteString(fmt.Sprintf("%s:%v\n", k, overrideValue))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s:%v\n", k, cv.boolValues[k]))
	}

	return sb.String()
}

// MarshalJSON marshal result to json format
func (cv ConfigVals) MarshalJSON() ([]byte, error) {
	var result struct {
		Int64Values  map[string]int64  `json:"int_64_values"`
		BoolValues   map[string]bool   `json:"bool_values"`
		StringValues map[string]string `json:"string_values"`
	}
	result.Int64Values = make(map[string]int64)
	result.BoolValues = make(map[string]bool)
	result.StringValues = make(map[string]string)

	// get and sort all keys including overrides
	int64Keys := make([]ConfigName, 0, len(cv.int64values)+len(int64Overrides))
	for k := range cv.int64values {
		int64Keys = append(int64Keys, k)
	}
	for k := range int64Overrides {
		if _, exists := cv.int64values[k]; !exists {
			int64Keys = append(int64Keys, k)
		}
	}
	sort.Slice(int64Keys, func(i, j int) bool {
		return int64Keys[i].String() < int64Keys[j].String()
	})

	// Same for bool and string keys
	boolKeys := make([]ConfigName, 0, len(cv.boolValues)+len(boolOverrides))
	for k := range cv.boolValues {
		boolKeys = append(boolKeys, k)
	}
	for k := range boolOverrides {
		if _, exists := cv.boolValues[k]; !exists {
			boolKeys = append(boolKeys, k)
		}
	}
	sort.Slice(boolKeys, func(i, j int) bool {
		return boolKeys[i].String() < boolKeys[j].String()
	})

	stringKeys := make([]ConfigName, 0, len(cv.stringValues)+len(stringOverrides))
	for k := range cv.stringValues {
		stringKeys = append(stringKeys, k)
	}
	for k := range stringOverrides {
		if _, exists := cv.stringValues[k]; !exists {
			stringKeys = append(stringKeys, k)
		}
	}
	sort.Slice(stringKeys, func(i, j int) bool {
		return stringKeys[i].String() < stringKeys[j].String()
	})

	for _, k := range int64Keys {
		if override, ok := int64Overrides[k]; ok {
			result.Int64Values[k.String()] = override
		} else {
			result.Int64Values[k.String()] = cv.int64values[k]
		}
	}

	for _, k := range boolKeys {
		if override, ok := boolOverrides[k]; ok {
			result.BoolValues[k.String()] = override
		} else {
			result.BoolValues[k.String()] = cv.boolValues[k]
		}
	}

	for _, k := range stringKeys {
		if override, ok := stringOverrides[k]; ok {
			result.StringValues[k.String()] = override
		} else {
			result.StringValues[k.String()] = cv.stringValues[k]
		}
	}

	return json.MarshalIndent(result, "", " ")
}

func GetSWVersion() (int64, error) {
	re := regexp.MustCompile(`\d+`) // matches digits
	versionParts := re.FindAllString(Version, -1)
	var version int64

	if len(versionParts) > 0 {
		majorVersion, err := strconv.ParseInt(versionParts[0], 10, 64)
		if err != nil {
			return 0, err
		}
		version = majorVersion
	} else {
		version = 1
	}
	return version, nil
}
