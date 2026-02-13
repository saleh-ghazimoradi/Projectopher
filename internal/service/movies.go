package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
)

type MovieService interface {
	CreateMovie(ctx context.Context, input *dto.CreateMovieReq) (*dto.MovieResp, error)
	GetMovie(ctx context.Context, id string) (*dto.MovieResp, error)
	GetMovies(ctx context.Context, page, limit int64) ([]dto.MovieResp, *helper.PaginatedMeta, error)
}

type movieService struct {
	movieRepository repository.MovieRepository
}

func (m *movieService) CreateMovie(ctx context.Context, input *dto.CreateMovieReq) (*dto.MovieResp, error) {
	movie := m.toMovie(input)
	if err := m.movieRepository.CreateMovie(ctx, movie); err != nil {
		return nil, err
	}
	return m.toMovieRepsDTO(movie), nil
}

func (m *movieService) GetMovie(ctx context.Context, id string) (*dto.MovieResp, error) {
	movie, err := m.movieRepository.GetMovie(ctx, id)
	if err != nil {
		return nil, err
	}
	return m.toMovieRepsDTO(movie), nil
}

func (m *movieService) GetMovies(ctx context.Context, page, limit int64) ([]dto.MovieResp, *helper.PaginatedMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	total, err := m.movieRepository.CountMovies(ctx)
	if err != nil {
		return nil, nil, err
	}

	movies, err := m.movieRepository.GetMovies(ctx, offset, limit)

	response := make([]dto.MovieResp, len(movies))
	for i := range movies {
		response[i] = *m.toMovieRepsDTO(&movies[i])
	}

	totalPages := (total + limit - 1) / limit
	meta := &helper.PaginatedMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPages,
	}

	return response, meta, nil
}

func (m *movieService) toMovie(input *dto.CreateMovieReq) *domain.Movie {
	genres := make([]domain.Genre, len(input.Genre))
	for i := range genres {
		genres[i] = domain.Genre{
			GenreId:   input.Genre[i].GenreId,
			GenreName: input.Genre[i].GenreName,
		}
	}

	ranking := domain.Ranking{
		RankingValue: input.Ranking.RankingValue,
		RankingName:  input.Ranking.RankingName,
	}
	return &domain.Movie{
		ImdbId:      input.ImdbId,
		Title:       input.Title,
		PosterPath:  input.PosterPath,
		YoutubeId:   input.YoutubeId,
		Genre:       genres,
		AdminReview: input.AdminReview,
		Ranking:     ranking,
	}
}

func (m *movieService) toMovieRepsDTO(movie *domain.Movie) *dto.MovieResp {
	genres := make([]dto.Genre, len(movie.Genre))
	for i := range genres {
		genres[i] = dto.Genre{
			GenreId:   movie.Genre[i].GenreId,
			GenreName: movie.Genre[i].GenreName,
		}
	}

	ranking := dto.Ranking{
		RankingValue: movie.Ranking.RankingValue,
		RankingName:  movie.Ranking.RankingName,
	}

	return &dto.MovieResp{
		Id:          movie.Id.Hex(),
		ImdbId:      movie.ImdbId,
		Title:       movie.Title,
		PosterPath:  movie.PosterPath,
		YoutubeId:   movie.YoutubeId,
		Genre:       genres,
		AdminReview: movie.AdminReview,
		Ranking:     ranking,
	}
}

func NewMovieService(movieRepository repository.MovieRepository) MovieService {
	return &movieService{
		movieRepository: movieRepository,
	}
}
