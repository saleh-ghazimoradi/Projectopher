package domain

type Movie struct {
	Id          string
	ImdbId      string
	Title       string
	PosterPath  string
	YoutubeId   string
	Genre       []Genre
	AdminReview string
	Ranking     Ranking
}
