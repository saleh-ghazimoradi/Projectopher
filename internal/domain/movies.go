package domain

import "go.mongodb.org/mongo-driver/v2/bson"

type Movie struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	ImdbId      string        `bson:"imdb_id"`
	Title       string        `bson:"title"`
	PosterPath  string        `bson:"poster_path"`
	YoutubeId   string        `bson:"youtube_id"`
	Genre       []Genre       `bson:"genre"`
	AdminReview string        `bson:"admin_review"`
	Ranking     Ranking       `bson:"ranking"`
}
