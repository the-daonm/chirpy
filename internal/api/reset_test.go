package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResetHandler_Forbidden(t *testing.T) {
	cfg := &ApiConfig{
		Platform: "prod",
	}
	
	req := httptest.NewRequest("POST", "/admin/reset", nil)
	rr := httptest.NewRecorder()

	cfg.ResetHandler(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Errorf("expected status forbidden, got %d", rr.Code)
	}
}
