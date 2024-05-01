package api

import "github.com/nahtann/go-lab/internal/api/router"

type ApiRouter interface {
	SetModules(modules []router.ApiRouterModule)
	SetRoute(
		method string,
		path string,
		handler router.FuncHandler,
	)
}
