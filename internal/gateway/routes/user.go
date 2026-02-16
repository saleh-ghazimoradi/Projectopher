package routes

import (
	"github.com/julienschmidt/httprouter"
	"github.com/saleh-ghazimoradi/Projectopher/internal/gateway/handlers"
	"net/http"
)

type UserRoute struct {
	userHandler *handlers.UserHandler
}

func (u *UserRoute) UserRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", u.userHandler.GetProfile)
	router.HandlerFunc(http.MethodGet, "/v1/users", u.userHandler.GetProfiles)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", u.userHandler.UpdateProfile)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", u.userHandler.DeleteProfile)
}

func NewUserRoute(userHandler *handlers.UserHandler) *UserRoute {
	return &UserRoute{
		userHandler: userHandler,
	}
}
