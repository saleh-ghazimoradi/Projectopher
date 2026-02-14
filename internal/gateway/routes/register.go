package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/middlewares"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"net/http"
)

type Register struct {
	healthRoute *HealthRoute
	authRoute   *AuthRoute
	movieRoute  *MovieRoute
	genreRoute  *GenreRoute
	rankRoute   *RankRoute
	userRoute   *UserRoute
	middlewares *middlewares.Middleware
}

type Options func(*Register)

func WithHealthRoute(healthRoute *HealthRoute) Options {
	return func(r *Register) {
		r.healthRoute = healthRoute
	}
}

func WithAuthRoute(authRoute *AuthRoute) Options {
	return func(r *Register) {
		r.authRoute = authRoute
	}
}

func WithMovieRoute(movieRoute *MovieRoute) Options {
	return func(r *Register) {
		r.movieRoute = movieRoute
	}
}

func WithGenreRoute(genreRoute *GenreRoute) Options {
	return func(r *Register) {
		r.genreRoute = genreRoute
	}
}

func WithRankRoute(rankRoute *RankRoute) Options {
	return func(r *Register) {
		r.rankRoute = rankRoute
	}
}

func WithUserRoute(userRoute *UserRoute) Options {
	return func(r *Register) {
		r.userRoute = userRoute
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
	r.authRoute.AuthRoutes(router)
	r.movieRoute.MovieRoutes(router)
	r.genreRoute.GenreRoutes(router)
	r.rankRoute.RankRoutes(router)
	r.userRoute.UserRoutes(router)
	return r.middlewares.Recover(r.middlewares.Logging(r.middlewares.CORS(r.middlewares.RateLimit(router))))
}

func NewRegister(opts ...Options) *Register {
	r := &Register{}
	for _, f := range opts {
		f(r)
	}
	return r
}
