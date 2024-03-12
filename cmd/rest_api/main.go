package restApi

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/nahtann/go-authentication/internal/api/modules"
	"github.com/nahtann/go-authentication/internal/api/router"
)

func RunServer() {
	// load .env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file")
		os.Exit(1)
	}

	// connect to the database
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	mux := http.NewServeMux()

	apiRouter := router.NewApiRouter(mux, "/api")

	modules := []router.ApiRouterModule{modules.AuthModule}
	apiRouter.SetModules(modules)

	log.Fatal(http.ListenAndServe(":3333", mux))
}
