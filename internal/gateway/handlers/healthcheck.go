package handlers

import (
	"github.com/saleh-ghazimoradi/Projectopher/config"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"net/http"
)

type HealthHandler struct {
	config *config.Config
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	env := map[string]any{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.config.Application.Environment,
			"version":     h.config.Application.Version,
		},
	}
	helper.SuccessResponse(w, "I'm breathing", env)
}

func NewHealthHandler(config *config.Config) *HealthHandler {
	return &HealthHandler{
		config: config,
	}
}
