package service

import (
	"context"
	"github.com/saleh-ghazimoradi/Projectopher/internal/domain"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
)

type UserService interface {
	GetProfile(ctx context.Context, id string) (*dto.UserResp, error)
	GetProfiles(ctx context.Context, page, limit int64) ([]dto.UserResp, *helper.PaginatedMeta, error)
	UpdateProfile(ctx context.Context, id string, input *dto.UpdateUserReq) (*dto.UserResp, error)
	DeleteProfile(ctx context.Context, id string) error
}

type userService struct {
	userRepository repository.UserRepository
}

func (u *userService) GetProfile(ctx context.Context, id string) (*dto.UserResp, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return u.toUser(user), nil
}

func (u *userService) GetProfiles(ctx context.Context, page, limit int64) ([]dto.UserResp, *helper.PaginatedMeta, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	total, err := u.userRepository.CountUser(ctx)
	if err != nil {
		return nil, nil, err
	}

	users, err := u.userRepository.GetUsers(ctx, offset, limit)
	if err != nil {
		return nil, nil, err
	}

	response := make([]dto.UserResp, len(users))
	for i := range users {
		response[i] = *u.toUser(&users[i])
	}

	totalPages := (total + limit - 1) / limit

	meta := &helper.PaginatedMeta{
		Page:      page,
		Limit:     limit,
		Total:     total,
		TotalPage: totalPages,
	}
	return response, meta, nil
}

func (u *userService) UpdateProfile(ctx context.Context, id string, input *dto.UpdateUserReq) (*dto.UserResp, error) {
	user, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.FirstName != nil {
		user.FirstName = *input.FirstName
	}

	if input.LastName != nil {
		user.LastName = *input.LastName
	}

	if err := u.userRepository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	return u.toUser(user), nil
}

func (u *userService) DeleteProfile(ctx context.Context, id string) error {
	return u.userRepository.DeleteUser(ctx, id)
}

func (u *userService) toUser(user *domain.User) *dto.UserResp {
	genres := make([]dto.Genre, len(user.FavoriteGenres))
	for i := range genres {
		genres[i] = dto.Genre{
			GenreId:   user.FavoriteGenres[i].GenreId,
			GenreName: user.FavoriteGenres[i].GenreName,
		}
	}

	return &dto.UserResp{
		Id:             user.Id,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Email:          user.Email,
		Role:           string(user.Role),
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
		FavoriteGenres: genres,
	}
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}
