package restApi

import (
	"log"
	"net/http"

	"github.com/nahtann/go-authentication/internal/api"
)

func RunServer() {
	mux := http.NewServeMux()

	api.SetRoutes(mux)

	log.Fatal(http.ListenAndServe(":3333", mux))
}
