package modules

import (
	"github.com/nahtann/go-authentication/internal/api/router"
	"github.com/nahtann/go-authentication/internal/handlers/auth_handlers"
	"github.com/nahtann/go-authentication/internal/utils"
)

const authRootRoute = "/auth"

func AuthModule(router *router.ApiRouter) {
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/signin"),
		auth_handlers.NewSignInHttpHandler(router.DB).Serve,
	)
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/signup"),
		auth_handlers.NewSignUpHttpHandler(router.DB).Serve,
	)
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/refresh-token"),
		auth_handlers.NewRefreshTokenHttpHandler(router.DB).Serve,
	)
}
