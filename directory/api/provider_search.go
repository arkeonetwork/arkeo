package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/arkeonetwork/arkeo/directory/utils"
)

// swagger:route Get /search searchProviders
//
// queries the service for a list of providers
//
// Parameters:
//   + name: chain
//     in: query
//     description: chain provider services
//     required: false
//     schema:
//      type: string
//   + name: pubkey
//     in: query
//     description: pubkey of provider
//     required: false
//     schema:
//      type: string
//   + name: sort
//     in: query
//     description: defines how to sort the list of providers
//     required: false
//     schema:
//      type: string
//      enum: age, conract_count, amount_paid
//   + name: max-distance
//     in: query
//     description: maximum distance in kilometers from provided coordinates
//     required: false
//     type: integer
//   + name: coordinates
//	   description: latitude and longitude (required when providing distance filter, example 40.7127837,-74.0059413)
//     in: query
//     required: false
//     type: string
//   + name: min-validator-payments
//	   description: minimum amount the provider has paid to validators
//     in: query
//     required: false
//	   type: integer
//   + name: min-provider-age
//	   description: minimum age of provider
//     in: query
//     required: false
//     type: integer
//   + name: min-free-rate-limit
//	   description: min rate limit for free tier of provider in requests per seconds
//     in: query
//     required: false
//	   type: integer
//   + name: min-payasyougo-rate-limit
//	   description: min rate limit for pay-as-you-go tier of provider in requests per seconds
//     in: query
//     required: false
//	   type: integer
//   + name: min-subscription-rate-limit
//	   description: min rate limit for subscription tier of provider in requests per seconds
//     in: query
//     required: false
//	   type: integer
//   + name: min-open-contracts
//	   description: minimum number of contracts open with proivder
//     in: query
//     required: false
//	   type: integer
// Responses:
//
//	200: ArkeoProviders
//	500: InternalServerError

func (a *ApiService) searchProviders(response http.ResponseWriter, request *http.Request) {

	sort := request.FormValue("sort")
	chain := request.FormValue("chain")
	pubkey := request.FormValue("pubkey")
	maxDistanceInput := request.FormValue("max-distance")
	coordinatesInput := request.FormValue("coordinates")
	minValidatorPaymentsInput := request.FormValue("min-validator-payments")
	minProviderAgeInput := request.FormValue("min-provider-age")
	minFreeRateLimitInput := request.FormValue("min-free-rate-limit")
	minPaygoRateLimitInput := request.FormValue("min-payasyougo-rate-limit")
	minSubscribeRateLimitInput := request.FormValue("min-subscription-rate-limit")
	minOpenContractsInput := request.FormValue("min-open-contracts")

	if (maxDistanceInput != "" && coordinatesInput == "") || (coordinatesInput != "" && maxDistanceInput == "") {
		respondWithError(response, http.StatusBadRequest, "max distance must accompany coordinates when supplied")
		return
	}

	searchParams := types.ProviderSearchParams{}

	switch sort {
	case string(types.ProviderSortKeyNone):
	case string(types.ProviderSortKeyAge):
		searchParams.SortKey = types.ProviderSortKeyAge
	case string(types.ProviderSortKeyAmountPaid):
		searchParams.SortKey = types.ProviderSortKeyAmountPaid
	case string(types.ProviderSortKeyContractCount):
		searchParams.SortKey = types.ProviderSortKeyContractCount
	default:
		respondWithError(response, http.StatusBadRequest, "sort key can not be parsed")
		return
	}

	searchParams.Pubkey = pubkey

	if chain != "" && !utils.ValidateChain(chain) {
		respondWithError(response, http.StatusBadRequest, fmt.Sprintf("%s is not a valid chain", chain))
	}
	searchParams.Chain = chain

	if maxDistanceInput != "" {
		var err error
		maxDistance, err := strconv.ParseInt(maxDistanceInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "max distance can not be parsed")
			return
		}
		coordinates, err := utils.ParseCoordinates(coordinatesInput)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "coordinates can not be parsed")
			return
		}
		searchParams.IsMaxDistanceSet = true
		searchParams.MaxDistance = maxDistance
		searchParams.Coordinates = coordinates
	}

	if minValidatorPaymentsInput != "" {
		var err error
		minValidatorPayments, err := strconv.ParseInt(minValidatorPaymentsInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-validator-payments can not be parsed")
			return
		}
		searchParams.IsMinValidatorPaymentsSet = true
		searchParams.MinValidatorPayments = minValidatorPayments
	}

	if minProviderAgeInput != "" {
		var err error
		minProviderAge, err := strconv.ParseInt(minProviderAgeInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-provider-age can not be parsed")
			return
		}
		searchParams.MinProviderAge = minProviderAge
		searchParams.IsMinProviderAgeSet = true
	}

	if minFreeRateLimitInput != "" {
		var err error
		minFreeRateLimit, err := strconv.ParseInt(minFreeRateLimitInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-free-rate-limit can not be parsed")
			return
		}
		searchParams.IsMinFreeRateLimitSet = true
		searchParams.MinFreeRateLimit = minFreeRateLimit
	}

	if minPaygoRateLimitInput != "" {
		var err error
		minPaygoRateLimit, err := strconv.ParseInt(minPaygoRateLimitInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-payasyougo-rate-limit can not be parsed")
			return
		}
		searchParams.IsMinPaygoRateLimitSet = true
		searchParams.MinPaygoRateLimit = minPaygoRateLimit
	}

	if minSubscribeRateLimitInput != "" {
		var err error
		minSubscribeRateLimit, err := strconv.ParseInt(minSubscribeRateLimitInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-subscription-rate-limit can not be parsed")
			return
		}
		searchParams.IsMinSubscribeRateLimitSet = true
		searchParams.MinSubscribeRateLimit = minSubscribeRateLimit
	}

	if minOpenContractsInput != "" {
		var err error
		minOpenContracts, err := strconv.ParseInt(minOpenContractsInput, 10, 64)
		if err != nil {
			respondWithError(response, http.StatusBadRequest, "min-open-contracts can not be parsed")
			return
		}
		searchParams.MinOpenContracts = minOpenContracts
		searchParams.IsMinOpenContractsSet = true
	}
	results, err := a.db.SearchProviders(searchParams)
	if err != nil {
		log.Errorf("error searching providers: %+v", err)
		respondWithError(response, http.StatusInternalServerError, "error searching providers")
	}

	respondWithJSON(response, http.StatusOK, results)
}
