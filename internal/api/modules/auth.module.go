package modules

import (
	"github.com/nahtann/go-authentication/internal/api/router"
	"github.com/nahtann/go-authentication/internal/handlers/auth_handlers"
	"github.com/nahtann/go-authentication/internal/storage/database/repositories"
)

const ROOT_ROUTE = "/auth"

func AuthModule(router *router.ApiRouter) {
	userRepository := repositories.NewUserRepository(router.DB)

	router.SetRoute(
		"POST",
		setSubRoute("/signin"),
		auth_handlers.NewSignInHttpHandler(userRepository).Serve,
	)
	router.SetRoute(
		"POST",
		setSubRoute("/signup"),
		auth_handlers.NewSignUpHttpHandler(userRepository).Serve,
	)
}

func setSubRoute(route string) string {
	return ROOT_ROUTE + route
}
