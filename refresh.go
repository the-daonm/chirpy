package main

import (
	"net/http"
	"time"

	"chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect bearer token")
		return
	}

	refreshToken, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect bearer token")
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "refresh token expired")
		return
	}
	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "refresh token revoked")
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	type refreshResponse struct {
		Token string `json:"token"`
	}
	responseToken := refreshResponse{
		Token: accessToken,
	}

	respondWithJSON(w, http.StatusOK, responseToken)
}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect bearer token")
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect bearer token")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
