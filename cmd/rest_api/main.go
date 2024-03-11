package restApi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nahtann/go-authentication/internal/api/router"
)

func RunServer() {
	mux := http.NewServeMux()

	apiRouter := router.NewApiRouter(mux, "/api")

	modules := []router.ApiRouterModule{auth}
	apiRouter.SetModules(modules)

	apiRouter.SetRoute("GET", "/", home)

	log.Fatal(http.ListenAndServe(":3333", mux))
}

func auth(router *router.ApiRouter) {
	router.SetRoute("POST", "/signin", signin)
}

func signin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Signin")
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home")
}
