package domain

import "time"

type Movie struct {
	Id          string
	ImdbId      string
	Title       string
	PosterPath  string
	YoutubeId   string
	Genres      []Genre
	AdminReview string
	Ranking     Ranking
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
