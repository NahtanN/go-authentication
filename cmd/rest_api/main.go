package rest_api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/nahtann/go-lab/internal/api/modules"
	"github.com/nahtann/go-lab/internal/api/router"
)

var (
	PORT      = ":3333"
	ROOT_PATH = "/api"
)

func RunServer() {
	// load .env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading .env file")
		os.Exit(1)
	}

	// database cononection pool
	dbpool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	mux := http.NewServeMux()

	apiRouter := router.NewApiRouter(mux, ROOT_PATH, dbpool)

	modules := []router.ApiRouterModule{modules.AuthModule, modules.UsersModule}
	apiRouter.SetModules(modules)

	fmt.Printf("Server start on: http://localhost%s%s", PORT, ROOT_PATH)
	log.Fatal(http.ListenAndServe(PORT, mux))
}
