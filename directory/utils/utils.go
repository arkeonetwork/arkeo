package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/arkeonetwork/arkeo/directory/sentinel"
	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/pkg/errors"

	"net/http"
	"net/url"

	resty "github.com/go-resty/resty/v2"
)

func ParseCoordinates(coordinates string) (types.Coordinates, error) {
	if coordinates == "" {
		return types.Coordinates{}, errors.New("empty string cannot be parsed into coordinates")
	}
	coordinatesSplit := strings.Split(coordinates, ",")
	if len(coordinatesSplit) != 2 {
		return types.Coordinates{}, errors.New("invalid string passed to coordinates")
	}
	latitude, err := strconv.ParseFloat(coordinatesSplit[0], 32)
	if err != nil {
		return types.Coordinates{}, errors.New("latitude cannot be parsed")
	}
	longitude, err := strconv.ParseFloat(coordinatesSplit[1], 32)
	if err != nil {
		return types.Coordinates{}, errors.New("longitude cannot be parsed")
	}
	return types.Coordinates{Latitude: latitude, Longitude: longitude}, nil
}

func ParseContractType(contractType string) (types.ContractType, error) {
	if types.ContractType(contractType) == types.ContractTypePayAsYouGo {
		return types.ContractType(contractType), nil
	} else if types.ContractType(contractType) == types.ContractTypeSubscription {
		return types.ContractType(contractType), nil
	} else {
		return types.ContractTypePayAsYouGo, fmt.Errorf("unexpected contract type %s", contractType)
	}
}

func IsNearEqual(a float64, b float64, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

// see arkeo-protocol/common/chain.go
var validChains = map[string]struct{}{"arkeo-mainnet-fullnode": {}, "btc-mainnet-fullnode": {}, "eth-mainnet-fullnode": {}, "gaia-mainnet-rpc-archive": {}, "swapi.dev": {}}

func ValidateChain(chain string) (ok bool) {
	_, ok = validChains[chain]
	return
}

func readFromNetwork(u *url.URL, retries int, maxBytes int) ([]byte, error) {
	client := resty.New()

	client.SetRetryCount(retries)
	client.SetTimeout(time.Second * 5)
	client.SetHeader("Accept", "application/json")
	resp, err := client.R().ForceContentType("application/json").Get(u.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("http status %d", resp.StatusCode())
	}

	body := resp.Body()
	if len(body) > maxBytes {
		return nil, errors.New("DownloadProviderMetadata: max bytes exceeded")
	}
	return body, nil
}

func readFromFilesystem(u *url.URL) (raw []byte, err error) {
	full := fmt.Sprintf("/%s%s", u.Host, u.Path)
	if raw, err = os.ReadFile(full); err != nil {
		return nil, errors.Wrapf(err, "error reading file %s", full)
	}
	return raw, nil
}

func DownloadProviderMetadata(metadataUrl string, retries int, maxBytes int) (*sentinel.Metadata, error) {

	u, err := url.Parse(metadataUrl)
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing url %s", metadataUrl)
	}

	var raw []byte
	switch u.Scheme {
	case "file":
		if raw, err = readFromFilesystem(u); err != nil {
			return nil, errors.Wrapf(err, "error reading metadata from fs")
		}
	default:
		if raw, err = readFromNetwork(u, retries, maxBytes); err != nil {
			return nil, errors.Wrapf(err, "error reading metadata from fs")
		}
	}

	result := &sentinel.Metadata{}
	if err = json.Unmarshal(raw, result); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling")
	}

	return result, nil
}
