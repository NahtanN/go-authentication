package rest_api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/nahtann/go-lab/docs"
	"github.com/nahtann/go-lab/internal/api/modules"
	"github.com/nahtann/go-lab/internal/api/router"
)

func RunServer(port, rootPath string) {
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

	mux.HandleFunc("GET /swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:3333/swagger/doc.json"),
	))

	apiRouter := router.NewApiRouter(mux, rootPath, dbpool)

	modules := []router.ApiRouterModule{modules.AuthModule, modules.UsersModule}
	apiRouter.SetModules(modules)

	fmt.Printf("Server start on: http://localhost%s%s", port, rootPath)
	log.Fatal(http.ListenAndServe(port, mux))
}
