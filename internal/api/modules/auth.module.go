package modules

import (
	"github.com/nahtann/go-authentication/internal/api/router"
	"github.com/nahtann/go-authentication/internal/handlers/auth_handlers"
)

const ROOT_ROUTE = "/auth"

func AuthModule(router *router.ApiRouter) {
	router.SetRoute("POST", setSubRoute("/signin"), auth_handlers.Signin)
	router.SetRoute(
		"POST",
		setSubRoute("/signup"),
		auth_handlers.NewSignUpHttpHandler(router.DB).Serve,
	)
}

func setSubRoute(route string) string {
	return ROOT_ROUTE + route
}
