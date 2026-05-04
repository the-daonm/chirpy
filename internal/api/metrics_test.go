package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewareMetricsInc(t *testing.T) {
	cfg := &ApiConfig{}
	
	handler := cfg.MiddlewareMetricsInc(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/app/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if cfg.FileserverHits.Load() != 1 {
		t.Errorf("expected 1 hit, got %d", cfg.FileserverHits.Load())
	}

	handler.ServeHTTP(rr, req)
	if cfg.FileserverHits.Load() != 2 {
		t.Errorf("expected 2 hits, got %d", cfg.FileserverHits.Load())
	}
}
