package modules

import (
	"github.com/nahtann/go-lab/internal/api/router"
	"github.com/nahtann/go-lab/internal/handlers/auth_handlers"
	"github.com/nahtann/go-lab/internal/utils"
)

const authRootRoute = "/auth"

func AuthModule(router *router.ApiRouter) {
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/signin"),
		auth_handlers.NewSignInHttpHandler(router.DB),
	)
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/signup"),
		auth_handlers.NewSignUpHttpHandler(router.DB),
	)
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/refresh-token"),
		auth_handlers.NewRefreshTokenHttpHandler(router.DB),
	)
}
