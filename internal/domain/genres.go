package domain

type Genre struct {
	GenreId   int    `bson:"genre_id"`
	GenreName string `bson:"genre_name"`
}
