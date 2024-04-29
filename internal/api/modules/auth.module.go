package modules

import (
	"github.com/nahtann/go-authentication/internal/api/router"
	"github.com/nahtann/go-authentication/internal/handlers/auth_handlers"
	"github.com/nahtann/go-authentication/internal/storage/database/repositories"
	"github.com/nahtann/go-authentication/internal/utils"
)

const authRootRoute = "/auth"

func AuthModule(router *router.ApiRouter) {
	userRepository := repositories.NewUserRepository(router.DB)
	refreshTokenRepository := repositories.NewRefreshTokenRepository(router.DB)

	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/signin"),
		auth_handlers.NewSignInHttpHandler(userRepository, refreshTokenRepository).Serve,
	)
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/signup"),
		auth_handlers.NewSignUpHttpHandler(userRepository).Serve,
	)
	router.SetRoute(
		"POST",
		utils.SetSubRoute(authRootRoute, "/refresh-token"),
		auth_handlers.NewRefreshTokenHttpHandler(router.DB).Serve,
	)
}
