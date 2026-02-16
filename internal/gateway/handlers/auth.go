package handlers

import (
	"errors"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
	"github.com/saleh-ghazimoradi/Projectopher/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

func (a *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var payload dto.RegisterReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid given payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRegister(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid request payload")
		return
	}

	user, err := a.authService.Register(r.Context(), &payload)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.BadRequestResponse(w, "Registration failed", err)
		case errors.Is(err, repository.ErrDuplicateEmail):
			helper.EditConflictResponse(w, "Registration failed", err)
		default:
			helper.InternalServerError(w, "Failed to register user", err)
		}
		return
	}
	helper.CreatedResponse(w, "user successfully registered", user)
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var payload dto.LoginReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid request payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateLogin(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid request payload")
		return
	}

	login, err := a.authService.Login(r.Context(), &payload)
	if err != nil {
		helper.InternalServerError(w, "Failed to login", err)
		return
	}

	helper.SuccessResponse(w, "Login successful", login)
}

func (a *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var payload dto.RefreshTokenReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid request payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRefreshToken(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid request payload")
		return
	}

	refreshToken, err := a.authService.RefreshToken(r.Context(), &payload)
	if err != nil {
		helper.InternalServerError(w, "Failed to refresh token", err)
		return
	}

	helper.SuccessResponse(w, "Refresh token successfully", refreshToken)
}

func (a *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var payload dto.RefreshTokenReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "Invalid request payload", err)
		return
	}

	v := helper.NewValidator()
	dto.ValidateRefreshToken(v, &payload)
	if !v.Valid() {
		helper.FailedValidationResponse(w, "Invalid request payload")
		return
	}

	if err := a.authService.Logout(r.Context(), &payload); err != nil {
		helper.InternalServerError(w, "Failed to logout", err)
		return
	}

	helper.SuccessResponse(w, "Logout successfully", nil)
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}
