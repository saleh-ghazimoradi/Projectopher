package service

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Projectopher/config"
	"github.com/saleh-ghazimoradi/Projectopher/infra/AI"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
)

type MovieService interface {
	CreateMovie(ctx context.Context, input *dto.CreateMovieReq) (*dto.MovieResp, error)
	GetMovie(ctx context.Context, id string) (*dto.MovieResp, error)
	GetMovies(ctx context.Context, page, limit int64) ([]dto.MovieResp, *helper.PaginatedMeta, error)
	UpdateAdminReview(ctx context.Context, imdbId string, input *dto.AdminReviewUpdateReq) (*dto.AdminReviewResp, error)
	GetRecommendedMovies(ctx context.Context, userId string) ([]dto.MovieResp, error)
	GetGenres(ctx context.Context) ([]dto.Genre, error)
}

type movieService struct {
	movieRepository   repository.MovieRepository
	rankingRepository repository.RankingRepository
	genreRepository   repository.GenreRepository
	userRepository    repository.UserRepository
	openAI            AI.OpenAI
	config            *config.Config
}

func (m *movieService) CreateMovie(ctx context.Context, input *dto.CreateMovieReq) (*dto.MovieResp, error) {
	movie := dto.FromCreateMovieReq(input)
	if err := m.movieRepository.CreateMovie(ctx, movie); err != nil {
		return nil, err
	}
	return dto.ToMovieResp(movie), nil
}

func (m *movieService) GetMovie(ctx context.Context, id string) (*dto.MovieResp, error) {
	movie, err := m.movieRepository.GetMovie(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.ToMovieResp(movie), nil
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
	for i, movie := range movies {
		response[i] = *dto.ToMovieResp(&movie)
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

func (m *movieService) UpdateAdminReview(ctx context.Context, imdbId string, input *dto.AdminReviewUpdateReq) (*dto.AdminReviewResp, error) {
	rankings, err := m.rankingRepository.GetRankings(ctx)
	if err != nil {
		return nil, err
	}

	var sentiments []string
	for _, r := range rankings {
		if r.RankingValue != 999 {
			sentiments = append(sentiments, r.RankingName)
		}
	}

	sentiment, err := m.openAI.GetSentiment(ctx, input.AdminReview, sentiments)
	if err != nil {
		return nil, err
	}

	var rankVal int
	for _, r := range rankings {
		if r.RankingName == sentiment {
			rankVal = r.RankingValue
			break
		}
	}

	if rankVal == 0 {
		return nil, errors.New("invalid sentiment ranking")
	}

	ranking := &domain.Ranking{
		RankingValue: rankVal,
		RankingName:  sentiment,
	}

	err = m.movieRepository.UpdateReview(ctx, imdbId, input.AdminReview, ranking)
	if err != nil {
		return nil, err
	}

	return &dto.AdminReviewResp{
		RankingName: sentiment,
		AdminReview: input.AdminReview,
	}, nil
}

func (m *movieService) GetRecommendedMovies(ctx context.Context, userId string) ([]dto.MovieResp, error) {
	limit := m.config.Application.MovieLimit
	if limit <= 0 {
		limit = 5
	}

	genres, err := m.userRepository.GetUserFavoriteGenres(ctx, userId)
	if err != nil {
		return nil, err
	}

	movies, err := m.movieRepository.GetRecommendedMovies(ctx, genres, limit)
	if err != nil {
		return nil, err
	}

	return dto.ToMoviesResp(movies), nil
}

func (m *movieService) GetGenres(ctx context.Context) ([]dto.Genre, error) {
	genres, err := m.genreRepository.GetGenres(ctx)
	if err != nil {
		return nil, err
	}
	return dto.ToGenresResp(genres), nil
}

func NewMovieService(movieRepository repository.MovieRepository, rankingRepository repository.RankingRepository, genreRepository repository.GenreRepository, userRepository repository.UserRepository, openAI AI.OpenAI, config *config.Config) MovieService {
	return &movieService{
		movieRepository:   movieRepository,
		rankingRepository: rankingRepository,
		genreRepository:   genreRepository,
		userRepository:    userRepository,
		openAI:            openAI,
		config:            config,
	}
}
