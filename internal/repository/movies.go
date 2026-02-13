package repository

import (
	"context"
	"errors"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
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
	if movie.Id.IsZero() {
		movie.Id = bson.NewObjectID()
	}

	if _, err := m.collection.InsertOne(ctx, movie); err != nil {
		return err
	}
	return nil
}

func (m *movieRepository) GetMovie(ctx context.Context, imdbId string) (*domain.Movie, error) {
	var movie domain.Movie
	if err := m.collection.FindOne(ctx, bson.M{"imdb_id": imdbId}).Decode(&movie); err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &movie, nil
}

func (m *movieRepository) GetMovies(ctx context.Context, offset, limit int64) ([]domain.Movie, error) {
	cursor, err := m.collection.Find(ctx, bson.M{}, options.Find().SetSkip(offset).SetLimit(limit))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var movies []domain.Movie
	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
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
