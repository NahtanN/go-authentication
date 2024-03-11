package modules

import (
	"github.com/nahtann/go-authentication/internal/api/router"
	"github.com/nahtann/go-authentication/internal/handlers/auth"
)

const ROOT_ROUTE = "/auth"

func AuthModule(router *router.ApiRouter) {
	router.SetRoute("POST", setSubRoute("/signin"), auth.Signin)
	router.SetRoute("POST", setSubRoute("/signup"), auth.Signup)
}

func setSubRoute(route string) string {
	return ROOT_ROUTE + route
}
