package utils

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/sentinel"

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

func ParseContractType(contractTypeStr string) (types.ContractType, error) {
	contractType := types.ContractType(contractTypeStr)
	switch contractType {
	case types.ContractTypePayAsYouGo:
	case types.ContractTypeSubscription:
	default:
		return contractType, fmt.Errorf("unexpected contract type %s", contractTypeStr)
	}
	return contractType, nil
}

func IsNearEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) <= epsilon
}

// see arkeo-protocol/common/service.go
var validServices = map[string]struct{}{"arkeo-mainnet-fullnode": {}, "btc-mainnet-fullnode": {}, "eth-mainnet-fullnode": {}, "gaia-mainnet-rpc-archive": {}, "mock": {}}

func ValidateService(service string) (ok bool) {
	_, ok = validServices[service]
	return
}

func readFromNetwork(u *url.URL, retries, maxBytes int) ([]byte, error) {
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

func DownloadProviderMetadata(metadataUrl string, retries, maxBytes int) (*sentinel.Metadata, error) {
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
			return nil, errors.Wrapf(err, "error reading metadata from network")
		}
	}

	// result := &sentinel.Metadata{}
	// if err = json.Unmarshal(raw, result); err != nil {
	// 	return nil, errors.Wrapf(err, "error unmarshaling")
	// }

	// Unmarshal into a generic map first
	var generic map[string]interface{}
	if err := json.Unmarshal(raw, &generic); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling to generic map")
	}

	// Fix service id types if needed
	if config, ok := generic["config"].(map[string]interface{}); ok {
		if services, ok := config["services"].([]interface{}); ok {
			for _, s := range services {
				if service, ok := s.(map[string]interface{}); ok {
					if idStr, isStr := service["id"].(string); isStr {
						if idInt, err := strconv.Atoi(idStr); err == nil {
							service["id"] = idInt
						}
					}
				}
			}
		}
	}

	// Marshal back to JSON
	fixedRaw, err := json.Marshal(generic)
	if err != nil {
		return nil, errors.Wrapf(err, "error re-marshaling fixed json")
	}

	// Now unmarshal into your strict struct
	result := &sentinel.Metadata{}
	if err := json.Unmarshal(fixedRaw, result); err != nil {
		return nil, errors.Wrapf(err, "error unmarshaling to Metadata struct")
	}

	return result, nil
}
