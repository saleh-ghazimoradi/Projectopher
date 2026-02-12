package routes

import (
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/handlers"
	"net/http"
)

type HealthRoute struct {
	healthHandler *handlers.HealthHandler
}

func (h *HealthRoute) HealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/healthcheck", h.healthHandler.HealthCheck)
}

func NewHealthRoute(healthHandler *handlers.HealthHandler) *HealthRoute {
	return &HealthRoute{
		healthHandler: healthHandler,
	}
}
