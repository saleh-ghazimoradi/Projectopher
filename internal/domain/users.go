package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type User struct {
	Id            bson.ObjectID `bson:"_id,omitempty"`
	UserId        string        `bson:"user_id"`
	FirstName     string        `bson:"first_name"`
	LastName      string        `bson:"last_name"`
	Email         string        `bson:"email"`
	Password      string        `bson:"password"`
	Role          string        `bson:"role"`
	CreatedAt     time.Time     `bson:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at"`
	Token         string        `bson:"token"`
	RefreshToken  string        `bson:"refresh_token"`
	FavoriteGenre []Genre       `bson:"favorite_genre"`
}
