package repository

import (
	"context"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
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

	var rankings []domain.Ranking
	if err := cursor.All(ctx, &rankings); err != nil {
		return nil, err
	}
	return rankings, nil
}

func NewRankingsRepository(database *mongo.Database, collectionName string) RankingRepository {
	return &rankingRepository{
		collection: database.Collection(collectionName),
	}
}
