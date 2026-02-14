package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type TokenRepository interface {
	CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error
	GetValidRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteRefreshTokenById(ctx context.Context, id string) error
	DeleteExpired(ctx context.Context) error
}

type tokenRepository struct {
	collection *mongo.Collection
}

func (t *tokenRepository) CreateRefreshToken(ctx context.Context, token *domain.RefreshToken) error {
	if token.Id.IsZero() {
		token.Id = bson.NewObjectID()
	}
	if _, err := t.collection.InsertOne(ctx, token); err != nil {
		return err
	}
	return nil
}

func (t *tokenRepository) GetValidRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var refreshToken domain.RefreshToken
	if err := t.collection.FindOne(ctx, bson.M{"token": token, "expires_at": bson.M{"$gt": time.Now()}}).Decode(&refreshToken); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &refreshToken, nil
}

func (t *tokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	if _, err := t.collection.DeleteOne(ctx, token); err != nil {
		return err
	}
	return nil
}

func (t *tokenRepository) DeleteRefreshTokenById(ctx context.Context, id string) error {
	if _, err := t.collection.DeleteOne(ctx, bson.M{"_id": id}); err != nil {
		return err
	}
	return nil
}

func (t *tokenRepository) DeleteExpired(ctx context.Context) error {
	if _, err := t.collection.DeleteMany(ctx, bson.M{"expires_at": bson.M{"$lte": time.Now()}}); err != nil {
		return err
	}
	return nil
}

func NewTokenRepository(database *mongo.Database, collectionName string) TokenRepository {
	return &tokenRepository{
		collection: database.Collection(collectionName),
	}
}
