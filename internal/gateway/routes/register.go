package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
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
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(helper.HTTPRouterNotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(helper.HTTPRouterMethodNotAllowedResponse)
	r.healthRoute.HealthRoutes(router)

	return r.middlewares.Recover(r.middlewares.Logging(r.middlewares.CORS(router)))
}

func NewRegister(opts ...Options) *Register {
	r := &Register{}
	for _, f := range opts {
		f(r)
	}
	return r
}
