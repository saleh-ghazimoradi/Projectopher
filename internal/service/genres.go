package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
)

type GenreService interface {
	GetGenres(ctx context.Context) ([]dto.Genre, error)
}

type genreService struct {
	genreRepository repository.GenreRepository
}

func (g *genreService) GetGenres(ctx context.Context) ([]dto.Genre, error) {
	genres, err := g.genreRepository.GetGenres(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.Genre, len(genres))
	for i := range genres {
		response[i] = *g.toGenres(&genres[i])
	}
	
	return response, nil
}

func (g *genreService) toGenres(genre *domain.Genre) *dto.Genre {
	return &dto.Genre{
		GenreId:   genre.GenreId,
		GenreName: genre.GenreName,
	}
}

func NewGenreService(genreRepository repository.GenreRepository) GenreService {
	return &genreService{
		genreRepository: genreRepository,
	}
}
