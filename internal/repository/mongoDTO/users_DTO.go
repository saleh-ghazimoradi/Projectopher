package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type UserDTO struct {
	Id             bson.ObjectID `bson:"_id,omitempty"`
	FirstName      string        `bson:"first_name"`
	LastName       string        `bson:"last_name"`
	Email          string        `bson:"email"`
	Password       string        `bson:"password"`
	Role           string        `bson:"role"`
	CreatedAt      time.Time     `bson:"created_at"`
	UpdatedAt      time.Time     `bson:"updated_at"`
	FavoriteGenres []GenreDTO    `bson:"favorite_genres"`
}

func FromUserCoreToDTO(input *domain.User) (*UserDTO, error) {
	var objectID bson.ObjectID
	var err error

	if input.Id != "" {
		objectID, err = bson.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, fmt.Errorf("invalid user id")
		}
	}

	genres := make([]GenreDTO, len(input.FavoriteGenres))
	for i, g := range input.FavoriteGenres {
		genres[i] = GenreDTO{
			GenreId:   g.GenreId,
			GenreName: g.GenreName,
		}
	}

	return &UserDTO{
		Id:             objectID,
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password:       input.Password,
		Role:           string(input.Role),
		CreatedAt:      input.CreatedAt,
		UpdatedAt:      input.UpdatedAt,
		FavoriteGenres: genres,
	}, nil
}

func FromUserDTOToCore(input *UserDTO) *domain.User {
	genres := make([]domain.Genre, len(input.FavoriteGenres))
	for i, g := range input.FavoriteGenres {
		genres[i] = domain.Genre{
			GenreId:   g.GenreId,
			GenreName: g.GenreName,
		}
	}

	return &domain.User{
		Id:             input.Id.Hex(),
		FirstName:      input.FirstName,
		LastName:       input.LastName,
		Email:          input.Email,
		Password:       input.Password,
		Role:           domain.UserRole(input.Role),
		CreatedAt:      input.CreatedAt,
		UpdatedAt:      input.UpdatedAt,
		FavoriteGenres: genres,
	}
}
