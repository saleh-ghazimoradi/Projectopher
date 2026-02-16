package repository

import (
	"context"
	"errors"
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

	var dtos []mongoDTO.MovieDTO
	if err := cursor.All(ctx, &dtos); err != nil {
		return nil, err
	}

	movies := make([]domain.Movie, len(dtos))
	for i := range dtos {
		movies[i] = *mongoDTO.FromMovieDTOToCore(&dtos[i])
	}

	return movies, nil
}

func (m *movieRepository) CountMovies(ctx context.Context) (int64, error) {
	return m.collection.CountDocuments(ctx, bson.M{})
}

func NewMovieRepository(database *mongo.Database, collectionName string) MovieRepository {
	return &movieRepository{
		collection: database.Collection(collectionName),
	}
}
