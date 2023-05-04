package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/arkeonetwork/arkeo/directory/db"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// swagger:model ArkeoContract
// type ArkeoContract2 struct {
// 	id uint64
// }

// swagger:model ArkeoContracts
type ArkeoContracts []*db.ArkeoContract

// swagger:route Get /contract/{pubkey} getContract
//
// Get a specific ArkeoContract by a unique id
//
// Parameters:
//   + name: id
//     in: path
//     description: contract id
//     required: true
//     type: string
//
// Responses:
//
//	200: ArkeoContract
//	500: InternalServerError

func (a *ApiService) getContract(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawId := vars["id"]
	id, err := strconv.ParseUint(rawId, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid contract id (%s)", rawId))
		return
	}

	contract, err := a.findContract(id)
	if err != nil {
		log.Errorf("error finding contract %d: %+v", id, err)
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("error finding contract with id %d", id))
		return
	}

	respondWithJSON(w, http.StatusOK, contract)
}

// find a contract by contract id
func (a *ApiService) findContract(id uint64) (*db.ArkeoContract, error) {
	dbContract, err := a.db.FindContract(id)
	if err != nil {
		return nil, errors.Wrapf(err, "error finding contract with id %d", id)
	}
	if dbContract == nil {
		return nil, nil
	}

	return dbContract, nil
}
