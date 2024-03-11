package auth

import (
	"fmt"
	"net/http"
	"os"
)

func Signup(w http.ResponseWriter, r *http.Request) error {
	secret := os.Getenv("JWT_SECRET")

	fmt.Println(secret)

	return nil
}
