package configs

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
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

var BlockTime = 5 * time.Second

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
	sb := strings.Builder{}
	// analyze-ignore(map-iteration)
	for k, v := range cv.int64values {
		if overrideValue, ok := int64Overrides[k]; ok {
			sb.WriteString(fmt.Sprintf("%s:%d\n", k, overrideValue))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s:%d\n", k, v))
	}
	// analyze-ignore(map-iteration)
	for k, v := range cv.boolValues {
		if overrideValue, ok := boolOverrides[k]; ok {
			sb.WriteString(fmt.Sprintf("%s:%v\n", k, overrideValue))
			continue
		}
		sb.WriteString(fmt.Sprintf("%s:%v\n", k, v))
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
	// analyze-ignore(map-iteration)
	for k, v := range cv.int64values {
		result.Int64Values[k.String()] = v
	}
	// analyze-ignore(map-iteration)
	for k, v := range int64Overrides {
		result.Int64Values[k.String()] = v
	}
	// analyze-ignore(map-iteration)
	for k, v := range cv.boolValues {
		result.BoolValues[k.String()] = v
	}
	// analyze-ignore(map-iteration)
	for k, v := range boolOverrides {
		result.BoolValues[k.String()] = v
	}
	// analyze-ignore(map-iteration)
	for k, v := range cv.stringValues {
		result.StringValues[k.String()] = v
	}
	// analyze-ignore(map-iteration)
	for k, v := range stringOverrides {
		result.StringValues[k.String()] = v
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
