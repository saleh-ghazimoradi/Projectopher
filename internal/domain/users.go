package domain

import (
	"time"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type User struct {
	Id             string
	FirstName      string
	LastName       string
	Email          string
	Password       string
	Role           UserRole
	CreatedAt      time.Time
	UpdatedAt      time.Time
	FavoriteGenres []Genre
}
