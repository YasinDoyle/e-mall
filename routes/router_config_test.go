package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPublicConfigRouteIsNotRegistered(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := NewRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/config/public", nil)

	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected public config route to be removed, got status %d", w.Code)
	}
}
