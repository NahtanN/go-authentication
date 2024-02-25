package main

import (
	"fmt"

	"github.com/joho/godotenv"

	restApi "github.com/nahtann/go-authentication/cmd/rest_api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	restApi.RunServer()
}
