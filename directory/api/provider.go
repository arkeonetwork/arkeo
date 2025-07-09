package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	"github.com/arkeonetwork/arkeo/directory/db"
)

// swagger:model ArkeoProvider
// type ArkeoProvider2 struct {
// 	Pubkey string
// }

// Contains info about a 500 Internal Server Error response
// swagger:model InternalServerError
type InternalServerError struct {
	Message string `json:"message"`
}

// swagger:model ArkeoProviders
type ArkeoProviders []*db.ArkeoProvider

// swagger:route Get /provider/{pubkey} getProvider
//
// Get a specific ArkeoProvider by a unique id (pubkey+service)
//
// Parameters:
//   + name: pubkey
//     in: path
//     description: provider public key
//     required: true
//     type: string
//   + name: service
//	   in: query
//     description: service identifier
//     required: true
//     type: string
//
// Responses:
//
//	200: ArkeoProvider
//	500: InternalServerError

func (a *ApiService) getProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pubkey := vars["pubkey"]
	service := r.FormValue("service")
	if pubkey == "" {
		respondWithError(w, http.StatusBadRequest, "pubkey is required")
		return
	}
	if service == "" {
		respondWithError(w, http.StatusBadRequest, "service is required")
		return
	}
	provider, err := a.findProvider(r.Context(), pubkey, service)
	if err != nil {
		log.Errorf("error finding provider for %s service %s: %+v", pubkey, service, err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error finding provider with pubkey %s", pubkey))
		return
	}

	respondWithJSON(w, http.StatusOK, provider)
}

// find a provider by pubkey+service
func (a *ApiService) findProvider(ctx context.Context, pubkey, service string) (*db.ArkeoProvider, error) {
	dbProvider, err := a.db.FindProvider(ctx, pubkey, service)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding provider for %s %s", pubkey, service)
	}
	if dbProvider == nil {
		return nil, nil
	}

	// provider := &db.ArkeoProvider{Pubkey: dbProvider.Pubkey}
	return dbProvider, nil
}

func (a *ApiService) getSubscriberContracts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pubkey := vars["pubkey"]
	service := r.FormValue("service")
	if pubkey == "" {
		respondWithError(w, http.StatusBadRequest, "pubkey is required")
		return
	}
	if service == "" {
		respondWithError(w, http.StatusBadRequest, "service is required")
		return
	}

	contracts, err := a.db.FindSubscriberContracts(r.Context(), pubkey, service)
	if err != nil {
		log.Errorf("error finding provider for %s service %s: %+v", pubkey, service, err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error finding subscriber with pubkey %s", pubkey))
		return
	}

	respondWithJSON(w, http.StatusOK, contracts)
}
