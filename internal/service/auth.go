package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/saleh-ghazimoradi/Projectopher/config"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
	"github.com/saleh-ghazimoradi/Projectopher/utils"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, input *dto.RegisterReq) (*dto.AuthResp, error)
	Login(ctx context.Context, input *dto.LoginReq) (*dto.AuthResp, error)
	RefreshToken(ctx context.Context, input *dto.RefreshTokenReq) (*dto.AuthResp, error)
	Logout(ctx context.Context, input *dto.RefreshTokenReq) error
}

type authService struct {
	config          *config.Config
	userRepository  repository.UserRepository
	tokenRepository repository.TokenRepository
}

func (a *authService) Register(ctx context.Context, input *dto.RegisterReq) (*dto.AuthResp, error) {
	existing, err := a.userRepository.GetUserByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, repository.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if existing != nil {
		return nil, repository.ErrDuplicateEmail
	}

	user := a.toUser(input)
	if err := a.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return a.generateAuthResp(ctx, user)
}

func (a *authService) Login(ctx context.Context, input *dto.LoginReq) (*dto.AuthResp, error) {
	user, err := a.userRepository.GetUserByEmail(ctx, input.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPassword(user.Password, input.Password) {
		return nil, errors.New("invalid credentials")
	}

	return a.generateAuthResp(ctx, user)
}

func (a *authService) RefreshToken(ctx context.Context, input *dto.RefreshTokenReq) (*dto.AuthResp, error) {
	claim, err := utils.ValidateToken(input.RefreshToken, a.config.JWT.Secret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	refreshToken, err := a.tokenRepository.GetValidRefreshToken(ctx, input.RefreshToken)
	if err != nil {
		return nil, errors.New("refresh token not found or expired")
	}

	user, err := a.userRepository.GetUserById(ctx, claim.UserId)
	if err != nil {
		return nil, err
	}

	if err := a.tokenRepository.DeleteRefreshTokenById(ctx, refreshToken.Id.Hex()); err != nil {
		return nil, err
	}

	return a.generateAuthResp(ctx, user)
}

func (a *authService) Logout(ctx context.Context, input *dto.RefreshTokenReq) error {
	return a.tokenRepository.DeleteRefreshTokenById(ctx, input.RefreshToken)
}

func (a *authService) toUser(input *dto.RegisterReq) *domain.User {
	favoriteGenres := make([]domain.Genre, len(input.FavoriteGenres))
	for i := range favoriteGenres {
		favoriteGenres[i] = domain.Genre{
			GenreId:   input.FavoriteGenres[i].GenreId,
			GenreName: input.FavoriteGenres[i].GenreName,
		}
	}
	hashedPassword, _ := utils.HashPassword(input.Password)
	return &domain.User{
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		Email:         input.Email,
		Password:      hashedPassword,
		Role:          domain.UserRoleUser,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		FavoriteGenre: favoriteGenres,
	}
}

func (a *authService) generateAuthResp(ctx context.Context, user *domain.User) (*dto.AuthResp, error) {
	accessToken, refreshToken, err := utils.GenerateToken(a.config, user.FirstName, user.LastName, user.Email, string(user.Role), user.Id.Hex())
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	refToken := &domain.RefreshToken{
		UserId:    user.Id,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(a.config.JWT.RefreshTokenExpires),
		CreatedAt: time.Now(),
	}

	if err := a.tokenRepository.CreateRefreshToken(ctx, refToken); err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	genres := make([]dto.Genre, len(user.FavoriteGenre))
	for i := range genres {
		genres[i] = dto.Genre{
			GenreId:   user.FavoriteGenre[i].GenreId,
			GenreName: user.FavoriteGenre[i].GenreName,
		}
	}

	return &dto.AuthResp{
		User: dto.UserResp{
			Id:            user.Id.Hex(),
			FirstName:     user.FirstName,
			LastName:      user.LastName,
			Email:         user.Email,
			Role:          string(user.Role),
			CreatedAt:     user.CreatedAt,
			UpdatedAt:     user.UpdatedAt,
			FavoriteGenre: genres,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func NewAuthService(config *config.Config, userRepository repository.UserRepository, tokenRepository repository.TokenRepository) AuthService {
	return &authService{
		config:          config,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
	}
}
