package handlers

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Projectopher/internal/dto"
	"github.com/saleh-ghazimoradi/Projectopher/internal/helper"
	"github.com/saleh-ghazimoradi/Projectopher/internal/repository"
	"github.com/saleh-ghazimoradi/Projectopher/internal/service"
	"net/http"
	"strconv"
)

type UserHandler struct {
	validator   *helper.Validator
	userService service.UserService
}

func (u *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		helper.BadRequestResponse(w, "Invalid id", errors.New("id is required"))
		return
	}

	profile, err := u.userService.GetProfile(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "Failed to fetch a movie")
		default:
			helper.InternalServerError(w, "Failed to fetch a movie", err)
		}
		return
	}

	helper.SuccessResponse(w, "profile successfully retrieved", profile)
}

func (u *UserHandler) GetProfiles(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
	if page < 0 {
		page = 1
	}

	limit, _ := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64)
	if limit < 0 {
		limit = 10
	}

	users, meta, err := u.userService.GetProfiles(r.Context(), page, limit)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "Failed to fetch a user")
		default:
			helper.InternalServerError(w, "Failed to fetch a user", err)
		}
		return
	}

	helper.PaginatedSuccessResponse(w, "Users successfully retrieved", users, *meta)
}

func (u *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		helper.BadRequestResponse(w, "Invalid id", errors.New("id is required"))
		return
	}

	var payload dto.UpdateUserReq
	if err := helper.ReadJSON(w, r, &payload); err != nil {
		helper.BadRequestResponse(w, "invalid payload", err)
		return
	}

	dto.ValidateUserUpdateReq(u.validator, &payload)
	if !u.validator.Valid() {
		helper.FailedValidationResponse(w, "Validation failed")
		return
	}

	updatedUser, err := u.userService.UpdateProfile(r.Context(), id, &payload)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "Failed to fetch a user")
		default:
			helper.InternalServerError(w, "Failed to fetch a user", err)
		}
		return
	}

	helper.SuccessResponse(w, "User successfully updated", updatedUser)
}

func (u *UserHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := httprouter.ParamsFromContext(r.Context()).ByName("id")
	if id == "" {
		helper.BadRequestResponse(w, "Invalid id", errors.New("id is required"))
		return
	}
	if err := u.userService.DeleteProfile(r.Context(), id); err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			helper.NotFoundResponse(w, "Failed to fetch a user")
		default:
			helper.InternalServerError(w, "Failed to fetch a user", err)
		}
		return
	}

	helper.SuccessResponse(w, "User successfully deleted", nil)
}

func NewUserHandler(validator *helper.Validator, userService service.UserService) *UserHandler {
	return &UserHandler{
		validator:   validator,
		userService: userService,
	}
}
