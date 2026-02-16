package dto

import (
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"time"
)

type UserResp struct {
	Id             string    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	FavoriteGenres []Genre   `json:"favorite_genres"`
}

type UpdateUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Password  string `json:"password"`
}

func ValidateUpdateUserReq(v *helper.Validator, req *UpdateUserReq) {
	validateUserNames(v, req.FirstName, req.LastName)
	validatePassword(v, req.Password)
}
