package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type UserDTO struct {
	Id            bson.ObjectID `bson:"_id,omitempty"`
	FirstName     string        `bson:"first_name"`
	LastName      string        `bson:"last_name"`
	Email         string        `bson:"email"`
	Password      string        `bson:"password"`
	Role          string        `bson:"role"`
	CreatedAt     time.Time     `bson:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at"`
	FavoriteGenre []GenreDTO    `bson:"favorite_genre"`
}

func FromUserCoreToDTO(input *domain.User) (*UserDTO, error) {
	id, err := bson.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %s", input.Id)
	}

	dto := &UserDTO{
		Id:            id,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		Password:      input.Password,
		Role:          string(input.Role),
		CreatedAt:     input.CreatedAt,
		UpdatedAt:     input.UpdatedAt,
		FavoriteGenre: make([]GenreDTO, len(input.FavoriteGenre)),
	}

	for i, g := range input.FavoriteGenre {
		dto.FavoriteGenre[i] = *FromGenreCoreToDTO(&g)
	}

	return dto, nil
}

func FromUserDTOToCore(input *UserDTO) *domain.User {
	core := &domain.User{
		Id:            input.Id.Hex(),
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		Password:      input.Password,
		Role:          domain.UserRole(input.Role),
		CreatedAt:     input.CreatedAt,
		UpdatedAt:     input.UpdatedAt,
		FavoriteGenre: make([]domain.Genre, len(input.FavoriteGenre)),
	}

	for i, g := range input.FavoriteGenre {
		core.FavoriteGenre[i] = *FromGenreDTOToCore(&g)
	}

	return core
}
