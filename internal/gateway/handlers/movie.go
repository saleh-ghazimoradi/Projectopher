package handlers

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
	"github.com/saleh-ghazimoradi/Projectopher/internal/service"
	"github.com/saleh-ghazimoradi/Projectopher/utils"
	"net/http"
	"strconv"
)

type MovieHandler struct {
	movieService service.MovieService
}

func (m *MovieHandler) AddMovie(w http.ResponseWriter, r *http.Request) {
	var payload dto.CreateMovieReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateCreateMovieReq(v, &payload)
	if !v.Valid() {
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

func (m *MovieHandler) AdminReviewUpdate(w http.ResponseWriter, r *http.Request) {
	role, exists := utils.RoleFromCtx(r.Context())
	if !exists {
		helper.BadRequestResponse(w, "Invalid role", nil)
		return
	}

	if role != string(domain.UserRoleAdmin) {
		helper.UnauthorizedResponse(w, "User must be part of the admin role")
		return
	}

	imdbId := r.URL.Query().Get("imdb_id")
	if imdbId == "" {
		helper.BadRequestResponse(w, "Invalid imdb_id", nil)
		return
	}

	var payload dto.AdminReviewUpdateReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateAdminReview(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Validation failed")
		return
	}

	resp, err := m.movieService.UpdateAdminReview(r.Context(), imdbId, &payload)
	if err != nil {
		helper.InternalServerError(w, "Failed to update admin review", err)
		return
	}

	helper.SuccessResponse(w, "Admin review successfully updated", resp)
}

func (m *MovieHandler) GetRecommendedMoviesHandler(w http.ResponseWriter, r *http.Request) {
	userId, exists := utils.UserIdFromCtx(r.Context())
	if !exists {
		helper.BadRequestResponse(w, "Invalid user id", nil)
		return
	}

	movies, err := m.movieService.GetRecommendedMovies(r.Context(), userId)
	if err != nil {
		helper.InternalServerError(w, "Failed to fetch recommended movies", err)
		return
	}

	helper.SuccessResponse(w, "Recommended movies successfully retrieved", movies)
}

func (m *MovieHandler) GetGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := m.movieService.GetGenres(r.Context())
	if err != nil {
		helper.InternalServerError(w, "Failed to fetch genres", err)
	}
	helper.SuccessResponse(w, "Genres successfully retrieved", genres)
}

func NewMovieHandler(movieService service.MovieService) *MovieHandler {
	return &MovieHandler{
		movieService: movieService,
	}
}
