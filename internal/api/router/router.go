package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FuncHandler interface {
	Serve(w http.ResponseWriter, r *http.Request) error
}

type MiddlewareFuncHandler interface {
	Serve(w http.ResponseWriter, r *http.Request) (*http.Request, bool, error)
}

type ApiRouterModule func(router *ApiRouter)

type ApiRouter struct {
	Mux      *http.ServeMux
	RootPath string
	modules  []ApiRouterModule
	DB       *pgxpool.Pool
}

func NewApiRouter(mux *http.ServeMux, rootPath string, dbStorage *pgxpool.Pool) *ApiRouter {
	return &ApiRouter{
		Mux:      mux,
		RootPath: rootPath,
		DB:       dbStorage,
	}
}

func (router *ApiRouter) SetModules(modules []ApiRouterModule) {
	router.modules = append(router.modules, modules...)

	router.build()
}

func (router *ApiRouter) SetRoute(
	method string,
	path string,
	handler FuncHandler,
	middlewares ...MiddlewareFuncHandler,
) {
	route := fmt.Sprintf("%s %s%s", strings.ToUpper(method), router.RootPath, path)

	router.Mux.HandleFunc(route, httpHandler(handler, middlewares...))
}

func (router *ApiRouter) build() {
	for _, fn := range router.modules {
		fn(router)
	}
}

func httpHandler(fn FuncHandler, middlewares ...MiddlewareFuncHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusInternalServerError

		req := r
		for _, middleware := range middlewares {
			newReq, next, err := middleware.Serve(w, req)
			if err != nil {
				http.Error(w, "Server Error", status)
			}

			if !next {
				return
			}

			req = newReq
		}

		err := fn.Serve(w, req)
		if err != nil {
			http.Error(w, "Server error", status)
		}
	}
}
