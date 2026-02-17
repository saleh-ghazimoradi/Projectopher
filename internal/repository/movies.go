package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository/mongoDTO"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MovieRepository interface {
	CreateMovie(ctx context.Context, movie *domain.Movie) error
	GetMovie(ctx context.Context, imdbId string) (*domain.Movie, error)
	GetMovies(ctx context.Context, offset, limit int64) ([]domain.Movie, error)
	GetRecommendedMovies(ctx context.Context, genres []string, limit int64) ([]domain.Movie, error)
	UpdateReview(ctx context.Context, imdbId string, adminReview string, ranking *domain.Ranking) error
	CountMovies(ctx context.Context) (int64, error)
}

type movieRepository struct {
	collection *mongo.Collection
}

func (m *movieRepository) CreateMovie(ctx context.Context, movie *domain.Movie) error {
	dto, err := mongoDTO.FromMovieCoreToDTO(movie)
	if err != nil {
		return err
	}

	result, err := m.collection.InsertOne(ctx, dto)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(bson.ObjectID); ok {
		movie.Id = oid.Hex()
	}

	return nil
}

func (m *movieRepository) GetMovie(ctx context.Context, imdbId string) (*domain.Movie, error) {
	var dto mongoDTO.MovieDTO

	err := m.collection.FindOne(ctx, bson.M{"imdb_id": imdbId}).Decode(&dto)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	return mongoDTO.FromMovieDTOToCore(&dto), nil
}

func (m *movieRepository) GetMovies(ctx context.Context, offset, limit int64) ([]domain.Movie, error) {
	cursor, err := m.collection.Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var DTOs []mongoDTO.MovieDTO
	if err := cursor.All(ctx, &DTOs); err != nil {
		return nil, err
	}

	movies := make([]domain.Movie, len(DTOs))
	for i := range DTOs {
		movies[i] = *mongoDTO.FromMovieDTOToCore(&DTOs[i])
	}

	return movies, nil
}

func (m *movieRepository) GetRecommendedMovies(ctx context.Context, genres []string, limit int64) ([]domain.Movie, error) {
	if len(genres) == 0 {
		return []domain.Movie{}, nil
	}

	filter := bson.D{
		{Key: "genre.genre_name", Value: bson.D{
			{Key: "$in", Value: genres},
		}},
	}

	opts := options.Find().
		SetSort(bson.D{{Key: "ranking.ranking_value", Value: 1}}).
		SetLimit(limit)

	cursor, err := m.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to find recommended movies: %w", err)
	}
	defer cursor.Close(ctx)

	var DTOs []mongoDTO.MovieDTO
	if err = cursor.All(ctx, &DTOs); err != nil {
		return nil, fmt.Errorf("failed to decode mongoDTO: %w", err)
	}

	movies := make([]domain.Movie, len(DTOs))
	for i := range DTOs {
		movies[i] = *mongoDTO.FromMovieDTOToCore(&DTOs[i])
	}

	return movies, nil
}

func (m *movieRepository) UpdateReview(ctx context.Context, imdbId string, adminReview string, ranking *domain.Ranking) error {
	_, err := m.collection.UpdateOne(ctx, bson.M{"imdb_id": imdbId}, bson.M{"$set": bson.M{
		"admin_review": adminReview,
		"ranking":      ranking,
	}})
	return err
}

func (m *movieRepository) CountMovies(ctx context.Context) (int64, error) {
	return m.collection.CountDocuments(ctx, bson.M{})
}

func NewMovieRepository(database *mongo.Database, collectionName string) MovieRepository {
	return &movieRepository{
		collection: database.Collection(collectionName),
	}
}
