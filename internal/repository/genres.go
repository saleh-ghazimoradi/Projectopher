package repository

import (
	"context"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type GenreRepository interface {
	GetGenres(ctx context.Context) ([]domain.Genre, error)
}

type genreRepository struct {
	collection *mongo.Collection
}

func (g *genreRepository) GetGenres(ctx context.Context) ([]domain.Genre, error) {
	cursor, err := g.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dtos []mongoDTO.GenreDTO
	if err = cursor.All(ctx, &dtos); err != nil {
		return nil, err
	}

	genres := make([]domain.Genre, len(dtos))
	for i := range dtos {
		genres[i] = *mongoDTO.FromGenreDTOToCore(&dtos[i])
	}

	return genres, nil
}

func NewGenresRepository(database *mongo.Database, collectionName string) GenreRepository {
	return &genreRepository{
		collection: database.Collection(collectionName),
	}
}
