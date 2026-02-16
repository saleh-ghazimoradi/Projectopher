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
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

func validateFirstName(v *helper.Validator, firstName *string) {
	if firstName != nil {
		v.Check(len(*firstName) >= 2, "firstname", "must be greater than two characters")
		v.Check(len(*firstName) <= 30, "firstname", "must not be greater than 30 characters")
	}
}

func validateLastName(v *helper.Validator, lastName *string) {
	if lastName != nil {
		v.Check(len(*lastName) >= 2, "lastname", "must be greater than two characters")
		v.Check(len(*lastName) <= 30, "lastname", "must not be greater than 30 characters")
	}
}

func ValidateUserUpdateReq(v *helper.Validator, req *UpdateUserReq) {
	validateFirstName(v, req.FirstName)
	validateLastName(v, req.LastName)
}
