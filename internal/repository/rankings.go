package repository

import (
	"context"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RankingRepository interface {
	GetRankings(ctx context.Context) ([]domain.Ranking, error)
}

type rankingRepository struct {
	collection *mongo.Collection
}

func (r *rankingRepository) GetRankings(ctx context.Context) ([]domain.Ranking, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var dtos []mongoDTO.RankingDTO
	if err = cursor.All(ctx, &dtos); err != nil {
		return nil, err
	}

	rankings := make([]domain.Ranking, len(dtos))
	for i := range dtos {
		rankings[i] = *mongoDTO.FromRankingDTOToCore(&dtos[i])
	}

	return rankings, nil
}

func NewRankingsRepository(database *mongo.Database, collectionName string) RankingRepository {
	return &rankingRepository{
		collection: database.Collection(collectionName),
	}
}
