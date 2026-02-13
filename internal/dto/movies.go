package dto

import "github.com/saleh-ghazimoradi/Projectopher/internal/helper"

type CreateMovieReq struct {
	ImdbId      string  `json:"imdb_id"`
	Title       string  `json:"title"`
	PosterPath  string  `json:"poster_path"`
	YoutubeId   string  `json:"youtube_id"`
	Genre       []Genre `json:"genre"`
	AdminReview string  `json:"admin_review"`
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
