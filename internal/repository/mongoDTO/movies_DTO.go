package mongoDTO

import (
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type MovieDTO struct {
	Id          bson.ObjectID `bson:"_id,omitempty"`
	ImdbId      string        `bson:"imdb_id"`
	Title       string        `bson:"title"`
	PosterPath  string        `bson:"poster_path"`
	YoutubeId   string        `bson:"youtube_id"`
	Genre       []GenreDTO    `bson:"genre"`
	AdminReview string        `bson:"admin_review"`
	Ranking     RankingDTO    `bson:"ranking"`
	CreatedAt   time.Time     `bson:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"`
}

func FromMovieCoreToDTO(input *domain.Movie) (*MovieDTO, error) {
	var oid bson.ObjectID
	var err error

	if input.Id != "" {
		oid, err = bson.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, err
		}
	}

	dto := &MovieDTO{
		Id:          oid,
		ImdbId:      input.ImdbId,
		Title:       input.Title,
		PosterPath:  input.PosterPath,
		YoutubeId:   input.YoutubeId,
		Genre:       make([]GenreDTO, len(input.Genres)),
		AdminReview: input.AdminReview,
		Ranking:     *FromRankingCoreToDTO(&input.Ranking),
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}

	for i := range input.Genres {
		dto.Genre[i] = *FromGenreCoreToDTO(&input.Genres[i])
	}

	return dto, nil
}

func FromMovieDTOToCore(input *MovieDTO) *domain.Movie {
	core := &domain.Movie{
		Id:          input.Id.Hex(),
		ImdbId:      input.ImdbId,
		Title:       input.Title,
		PosterPath:  input.PosterPath,
		YoutubeId:   input.YoutubeId,
		Genres:      make([]domain.Genre, len(input.Genre)),
		AdminReview: input.AdminReview,
		Ranking:     *FromRankingDTOToCore(&input.Ranking),
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}

	for i, g := range input.Genre {
		core.Genres[i] = *FromGenreDTOToCore(&g)
	}

	return core
}
