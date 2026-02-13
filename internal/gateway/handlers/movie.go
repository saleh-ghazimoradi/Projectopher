package handlers

import (
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
	"github.com/saleh-ghazimoradi/Projectopher/internal/service"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	movieService service.MovieService
	validator    *helper.Validator
}

func (m *MovieHandler) AddMovie(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateMovieReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	dto.ValidateCreateMovieReq(m.validator, &payload)
	if !m.validator.Valid() {
		helper.FailedValidationResponse(w, "Validation failed")
		return
	}

	movie, err := m.movieService.CreateMovie(r.Context(), &payload)
	if err != nil {
		helper.InternalServerError(w, "failed to add movie", err)
		return
	}

	helper.CreatedResponse(w, "Movie successfully added", movie)
}

func (m *MovieHandler) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("imdb_id")
	if id == "" {
		helper.BadRequestResponse(w, "Invalid id", errors.New("id is required"))
		return
	}

	fmt.Printf("id: %s, type of id: %T\n", id, id)

	movie, err := m.movieService.GetMovie(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "Movie not found")
		default:
			helper.InternalServerError(w, "Failed to fetch a movie", err)
		}
		return
	}

	helper.SuccessResponse(w, "Movie successfully fetched", movie)
}

func (m *MovieHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if page < 0 {
		page = 1
	}

	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit < 0 {
		limit = 10
	}

	movies, meta, err := m.movieService.GetMovies(r.Context(), page, limit)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "Movie not found")
		default:
			helper.InternalServerError(w, "Failed to fetch movies", err)
		}
		return
	}

	helper.PaginatedSuccessResponse(w, "Movies successfully retrieved", movies, *meta)
}

func NewMovieHandler(movieService service.MovieService, validator *helper.Validator) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
		validator:    validator,
	}
}
