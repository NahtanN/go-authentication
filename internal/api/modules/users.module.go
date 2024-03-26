package modules

import (
	"github.com/nahtann/go-authentication/internal/api/router"
	"github.com/nahtann/go-authentication/internal/handlers/users_handlers"
	"github.com/nahtann/go-authentication/internal/middlewares"
	"github.com/nahtann/go-authentication/internal/storage/database/repositories"
	"github.com/nahtann/go-authentication/internal/utils"
)

const usersRootRoute = "/users"

func UsersModule(router *router.ApiRouter) {
	userRepository := repositories.NewUserRepository(router.DB)

	router.SetRoute(
		"GET",
		utils.SetSubRoute(usersRootRoute, "/current"),
		users_handlers.NewCurrentUserHttpHandler(userRepository).Server,
		middlewares.NewJWTValidationHttpHandler().Serve,
	)
}
