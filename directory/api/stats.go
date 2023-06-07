package api

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/arkeonetwork/arkeo/directory/types"
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
	arkeoStats, err := a.db.GetArkeoNetworkStats(r.Context())
	if err != nil {
		log.Error("error finding stats for Arkeo Network")
		respondWithError(w, http.StatusInternalServerError, "error finding stats for Arkeo Network")
		return
	}
	respondWithJSON(w, http.StatusOK, arkeoStats)
}

// swagger:route Get /stats/{service} getStatsService
//
// get service specific network stats
// Parameters:
//   + name: service
//     in: path
//     description: service identifier
//     required: true
//     type: string
//
// Responses:
//
//	200: ServiceStats
//	500: InternalServerError

func getStatsService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	service := vars["service"]
	if service == "" {
		respondWithError(w, http.StatusBadRequest, "service is required")
		return
	}
	respondWithJSON(w, http.StatusOK, &types.ServiceStats{})
}
