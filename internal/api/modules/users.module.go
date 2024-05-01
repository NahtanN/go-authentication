package modules

import (
	"github.com/nahtann/go-lab/internal/api/router"
	"github.com/nahtann/go-lab/internal/handlers/users_handlers"
	"github.com/nahtann/go-lab/internal/middlewares"
	"github.com/nahtann/go-lab/internal/utils"
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
