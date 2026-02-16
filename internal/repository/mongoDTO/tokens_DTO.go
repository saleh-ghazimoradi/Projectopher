package mongoDTO

import (
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
	userOID, err := bson.ObjectIDFromHex(input.UserId)
	if err != nil {
		return nil, err
	}

	var tokenOID bson.ObjectID
	if input.Id != "" {
		tokenOID, err = bson.ObjectIDFromHex(input.Id)
		if err != nil {
			return nil, err
		}
	}

	return &RefreshTokenDTO{
		Id:        tokenOID,
		UserId:    userOID,
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
