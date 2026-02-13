package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	Id            bson.ObjectID `bson:"_id,omitempty"`
	UserId        string        `bson:"user_id"`
	FirstName     string        `bson:"first_name"`
	LastName      string        `bson:"last_name"`
	Email         string        `bson:"email"`
	Password      string        `bson:"password"`
	Role          UserRole      `bson:"role"`
	CreatedAt     time.Time     `bson:"created_at"`
	UpdatedAt     time.Time     `bson:"updated_at"`
	FavoriteGenre []Genre       `bson:"favorite_genre"`
}

type RefreshToken struct {
	Id        bson.ObjectID `bson:"_id,omitempty"`
	UserId    bson.ObjectID `bson:"user_id"`
	Token     string        `bson:"token"`
	ExpiresAt time.Time     `bson:"expires_at"`
	CreatedAt time.Time     `bson:"created_at"`
}
