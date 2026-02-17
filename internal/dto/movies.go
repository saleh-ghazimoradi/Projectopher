package dto

import (
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
)

type CreateMovieReq struct {
	ImdbId      string  `json:"imdb_id"`
	Title       string  `json:"title"`
	PosterPath  string  `json:"poster_path"`
	YoutubeId   string  `json:"youtube_id"`
	AdminReview string  `json:"admin_review"`
	Genre       []Genre `json:"genre"`
	Ranking     Ranking `json:"ranking"`
}

type MovieResp struct {
	Id          string  `json:"id"`
	ImdbId      string  `json:"imdb_id"`
	Title       string  `json:"title"`
	PosterPath  string  `json:"poster_path"`
	YoutubeId   string  `json:"youtube_id"`
	Genre       []Genre `json:"genre"`
	AdminReview string  `json:"admin_review"`
	Ranking     Ranking `json:"ranking"`
}

func ToMovieResp(movie *domain.Movie) *MovieResp {
	genres := make([]Genre, len(movie.Genres))
	for i, g := range movie.Genres {
		genres[i] = Genre{
			GenreId:   g.GenreId,
			GenreName: g.GenreName,
		}
	}

	return &MovieResp{
		Id:          movie.Id,
		ImdbId:      movie.ImdbId,
		Title:       movie.Title,
		PosterPath:  movie.PosterPath,
		YoutubeId:   movie.YoutubeId,
		Genre:       genres,
		AdminReview: movie.AdminReview,
		Ranking: Ranking{
			RankingValue: movie.Ranking.RankingValue,
			RankingName:  movie.Ranking.RankingName,
		},
	}
}

func ToMoviesResp(movies []domain.Movie) []MovieResp {
	DTOs := make([]MovieResp, len(movies))
	for i, movie := range movies {
		DTOs[i] = *ToMovieResp(&movie)
	}
	return DTOs
}

func FromCreateMovieReq(dto *CreateMovieReq) *domain.Movie {
	genres := make([]domain.Genre, len(dto.Genre))
	for i, g := range dto.Genre {
		genres[i] = domain.Genre{
			GenreId:   g.GenreId,
			GenreName: g.GenreName,
		}
	}
	return &domain.Movie{
		ImdbId:      dto.ImdbId,
		Title:       dto.Title,
		PosterPath:  dto.PosterPath,
		YoutubeId:   dto.YoutubeId,
		Genres:      genres,
		AdminReview: dto.AdminReview,
		Ranking: domain.Ranking{
			RankingValue: dto.Ranking.RankingValue,
			RankingName:  dto.Ranking.RankingName,
		},
	}
}

func ToGenreResp(genre *domain.Genre) *Genre {
	return &Genre{
		GenreId:   genre.GenreId,
		GenreName: genre.GenreName,
	}
}

func ToGenresResp(genres []domain.Genre) []Genre {
	DTOs := make([]Genre, len(genres))
	for i, g := range genres {
		DTOs[i] = *ToGenreResp(&g)
	}
	return DTOs
}

func validateImdbId(v *helper.Validator, imdbId string) {
	v.Check(imdbId != "", "imdbId", "imdbId is required")
}

func validateTitle(v *helper.Validator, title string) {
	v.Check(title != "", "title", "title must be provided")
	v.Check(len(title) > 2, "title", "title must be at least 2 characters")
	v.Check(len(title) < 500, "title", "title must not be greater 500 characters")
}

func validatePosterPath(v *helper.Validator, posterPath string) {
	v.Check(posterPath != "", "posterPath", "posterPath must be provided")
	v.Check(helper.IsURL(posterPath), "posterPath", "posterPath must be a URL")
}

func validateYoutubeId(v *helper.Validator, youtubeId string) {
	v.Check(youtubeId != "", "youtubeId", "youtubeId must be provided")
}

func validateGenre(v *helper.Validator, genre []Genre) {
	v.Check(genre != nil, "genre", "genre must be provided")
	v.Check(len(genre) >= 1, "genre", "must contain at least 1 genre")
	v.Check(len(genre) <= 5, "genre", "must not contain more than 5 genres")
	v.Check(helper.Unique(genre), "genres", "must not contain duplicate values")
}

func validateRanking(v *helper.Validator, ranking *Ranking) {
	v.Check(ranking != nil, "ranking", "ranking must be provided")
}

func ValidateCreateMovieReq(v *helper.Validator, req *CreateMovieReq) {
	validateImdbId(v, req.ImdbId)
	validateTitle(v, req.Title)
	validatePosterPath(v, req.PosterPath)
	validateYoutubeId(v, req.YoutubeId)
	validateGenre(v, req.Genre)
	validateRanking(v, &req.Ranking)
}
