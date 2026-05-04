package api

import (
	"sync/atomic"

	"chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
	Platform       string
	JwtSecret      string
	PolkaKey       string
}
