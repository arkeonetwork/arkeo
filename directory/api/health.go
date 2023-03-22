package api

import "net/http"

type Health struct {
	Overall string
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, Health{Overall: "LGTM"})
}
