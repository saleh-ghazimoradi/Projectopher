package routes

import (
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/middlewares"
	"net/http"
)

type Register struct {
	healthRoute *HealthRoute
	middlewares *middlewares.Middleware
}

type Options func(*Register)

func WithHealthRoute(healthRoute *HealthRoute) Options {
	return func(r *Register) {
		r.healthRoute = healthRoute
	}
}

func WithMiddleware(middlewares *middlewares.Middleware) Options {
	return func(r *Register) {
		r.middlewares = middlewares
	}
}

func (r *Register) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()
	r.healthRoute.HealthRoutes(mux)
	return r.middlewares.Recover(r.middlewares.Logging(r.middlewares.CORS(mux)))
}

func NewRegister(opts ...Options) *Register {
	r := &Register{}
	for _, f := range opts {
		f(r)
	}
	return r
}
