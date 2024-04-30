package modules

import (
	"github.com/nahtann/go-authentication/internal/api/router"
	"github.com/nahtann/go-authentication/internal/handlers/users_handlers"
	"github.com/nahtann/go-authentication/internal/middlewares"
	"github.com/nahtann/go-authentication/internal/utils"
)

const usersRootRoute = "/users"

func UsersModule(router *router.ApiRouter) {
	router.SetRoute(
		"GET",
		utils.SetSubRoute(usersRootRoute, "/current"),
		users_handlers.NewCurrentUserHttpHandler(router.DB),
		middlewares.NewJWTValidationHttpHandler(),
	)
}
