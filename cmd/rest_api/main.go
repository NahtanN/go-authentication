package restApi

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/nahtann/go-authentication/internal/api/modules"
	"github.com/nahtann/go-authentication/internal/api/router"
)

func RunServer() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	mux := http.NewServeMux()

	apiRouter := router.NewApiRouter(mux, "/api")

	modules := []router.ApiRouterModule{modules.AuthModule}
	apiRouter.SetModules(modules)

	log.Fatal(http.ListenAndServe(":3333", mux))
}
