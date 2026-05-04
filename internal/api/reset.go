package api

import (
	"net/http"
)

func (cfg *ApiConfig) ResetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Platform different from dev")
		return
	}

	err := cfg.DB.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not delete all users")
		return
	}

	w.WriteHeader(http.StatusOK)
}
