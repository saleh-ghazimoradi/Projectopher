package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/handlers"
	"net/http"
)

type MovieRoute struct {
	movieHandler *handlers.MovieHandler
}

func (m *MovieRoute) MovieRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/movies", m.movieHandler.AddMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:imdb_id", m.movieHandler.GetMovie)
	router.HandlerFunc(http.MethodGet, "/v1/movies", m.movieHandler.GetMovies)

}

func NewMovieRoute(movieHandler *handlers.MovieHandler) *MovieRoute {
	return &MovieRoute{
		movieHandler: movieHandler,
	}
}
