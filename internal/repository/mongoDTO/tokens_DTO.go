package mongoDTO

import (
	"fmt"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type RefreshTokenDTO struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	UserId    bson.ObjectID `bson:"user_id"`
	Token     string        `bson:"token"`
	ExpiresAt time.Time     `bson:"expires_at"`
	CreatedAt time.Time     `bson:"created_at"`
}

func FromRefreshTokenCoreToDTO(input *domain.RefreshToken) (*RefreshTokenDTO, error) {
	id, err := bson.ObjectIDFromHex(input.Id)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token id: %s", input.Id)
	}

	userId, err := bson.ObjectIDFromHex(input.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %s", input.UserId)
	}

	return &RefreshTokenDTO{
		Id:        id,
		UserId:    userId,
		Token:     input.Token,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: input.CreatedAt,
	}, nil
}

func FromRefreshTokenDTOToCore(input *RefreshTokenDTO) *domain.RefreshToken {
	return &domain.RefreshToken{
		Id:        input.Id.Hex(),
		UserId:    input.UserId.Hex(),
		Token:     input.Token,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: input.CreatedAt,
	}
}
