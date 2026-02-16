package mongoDTO

import (
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
)

type GenreDTO struct {
	GenreId   int    `bson:"genre_id"`
	GenreName string `bson:"genre_name"`
}

func FromGenreCoreToDTO(input *domain.Genre) *GenreDTO {
	return &GenreDTO{
		GenreId:   input.GenreId,
		GenreName: input.GenreName,
	}
}

func FromGenreDTOToCore(input *GenreDTO) *domain.Genre {
	return &domain.Genre{
		GenreId:   input.GenreId,
		GenreName: input.GenreName,
	}
}
