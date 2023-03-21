package api

import (
	"net/http"

	"github.com/arkeonetwork/arkeo/directory/types"
	"github.com/gorilla/mux"
)

// swagger:route Get /stats getStatsArkeo
//
// get Arkeo network stats
//
// Responses:
//
//	200: ArkeoStats
//	500: InternalServerError
func (a *ApiService) getStatsArkeo(w http.ResponseWriter, r *http.Request) {
	arkeoStats, err := a.db.GetArkeoNetworkStats()
	if err != nil {
		log.Error("error finding stats for Arkeo Network")
		respondWithError(w, http.StatusInternalServerError, "error finding stats for Arkeo Network")
		return
	}
	respondWithJSON(w, http.StatusOK, arkeoStats)
}

// swagger:route Get /stats/{chain} getStatsChain
//
// get chain specific network stats
// Parameters:
//   + name: chain
//     in: path
//     description: chain identifier
//     required: true
//     type: string
//
// Responses:
//
//	200: ChainStats
//	500: InternalServerError

func getStatsChain(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chain := vars["chain"]
	if chain == "" {
		respondWithError(w, http.StatusBadRequest, "chain is required")
		return
	}
	respondWithJSON(w, http.StatusOK, &types.ChainStats{})
}
