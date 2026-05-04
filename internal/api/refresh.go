package api

import (
	"net/http"
	"time"

	"chirpy/internal/auth"
)

func (cfg *ApiConfig) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	refreshToken, err := cfg.DB.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Refresh token expired")
		return
	}
	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Refresh token revoked")
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.JwtSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create access token")
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

func (cfg *ApiConfig) RevokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err = cfg.DB.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
