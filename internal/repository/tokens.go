package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository/mongoDTO"
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
	dto, err := mongoDTO.FromRefreshTokenCoreToDTO(token)
	if err != nil {
		return err
	}

	result, err := t.collection.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		token.Id = oid.Hex()
	}

	return nil
}

func (t *tokenRepository) GetValidRefreshToken(ctx context.Context, token string) (*domain.RefreshToken, error) {
	var dto mongoDTO.RefreshTokenDTO

	err := t.collection.FindOne(ctx, bson.M{
		"token":      token,
		"expires_at": bson.M{"$gt": time.Now()},
	}).Decode(&dto)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return mongoDTO.FromRefreshTokenDTOToCore(&dto), nil
}

func (t *tokenRepository) DeleteRefreshToken(ctx context.Context, token string) error {
	_, err := t.collection.DeleteOne(ctx, bson.M{"token": token})
	return err
}

func (t *tokenRepository) DeleteRefreshTokenById(ctx context.Context, id string) error {
	oid, _ := t.oId(id)
	_, err := t.collection.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

func (t *tokenRepository) DeleteExpired(ctx context.Context) error {
	_, err := t.collection.DeleteMany(ctx, bson.M{
		"expires_at": bson.M{"$lte": time.Now()},
	})
	return err
}

func (t *tokenRepository) oId(id string) (bson.ObjectID, error) {
	oid, err := bson.ObjectIDFromHex(id)
	return oid, err
}

func NewTokenRepository(database *mongo.Database, collectionName string) TokenRepository {
	return &tokenRepository{
		collection: database.Collection(collectionName),
	}
}
