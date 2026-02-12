package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/handlers"
	"net/http"
)

type HealthRoute struct {
	healthHandler *handlers.HealthHandler
}

func (h *HealthRoute) HealthRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", h.healthHandler.HealthCheck)
}

func NewHealthRoute(healthHandler *handlers.HealthHandler) *HealthRoute {
	return &HealthRoute{
		healthHandler: healthHandler,
	}
}
