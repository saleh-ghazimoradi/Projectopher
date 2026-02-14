package dto

import (
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"time"
)

type RegisterReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResp struct {
	Id            string    `json:"id"`
	UserId        string    `json:"user_id"`
	FirstName     string    `json:"first_name"`
	LastName      string    `json:"last_name"`
	Email         string    `json:"email"`
	Role          string    `json:"role"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	FavoriteGenre []Genre   `json:"favorite_genre"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type TokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func validateEmail(v *helper.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(helper.Matches(email, helper.EmailRX), "email", "must be a valid email address")
}

func validatePassword(v *helper.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func validateUserNames(v *helper.Validator, firstName string, lastName string) {
	v.Check(firstName != "", "first_name", "must be provided")
	v.Check(lastName != "", "last_name", "must be provided")
	v.Check(len(firstName) >= 2, "first_name", "must be at least 2 characters long")
	v.Check(len(lastName) >= 2, "last_name", "must be at least 2 characters long")
	v.Check(len(firstName) <= 72, "first_name", "must not be more than 72 bytes long")
	v.Check(len(lastName) <= 72, "last_name", "must not be more than 72 bytes long")
}

func validateRole(v *helper.Validator, role string) {
	v.Check(role != "", "role", "must be provided")
	v.Check(helper.PermittedValue(role, "user", "admin"), "role", "must be either user or admin")
}

func ValidateLogin(v *helper.Validator, req *LoginReq) {
	validateEmail(v, req.Email)
	validatePassword(v, req.Password)
}

func ValidateRegister(v *helper.Validator, req *RegisterReq) {
	validateUserNames(v, req.FirstName, req.LastName)
	validateEmail(v, req.Email)
	validatePassword(v, req.Password)
	validateRole(v, req.Role)
}
