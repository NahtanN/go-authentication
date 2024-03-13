package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FuncHandler func(w http.ResponseWriter, r *http.Request) error

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
) {
	route := fmt.Sprintf("%s %s%s", strings.ToUpper(method), router.RootPath, path)

	router.Mux.HandleFunc(route, httpHandler(handler))
}

func (router *ApiRouter) build() {
	for _, fn := range router.modules {
		fn(router)
	}
}

func httpHandler(fn FuncHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			status := http.StatusInternalServerError

			http.Error(w, "Server error", status)
		}
	}
}
