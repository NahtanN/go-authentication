package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

func SetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /public", public)
	mux.HandleFunc("POST /private", private)
	mux.HandleFunc("POST /login", login)
}

func public(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a public route")
}

func private(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var message Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		return
	}

	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "This is a private route")
}

func login(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("JWT_SECRET")
	byteSecret := []byte(secret)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(10 * time.Minute)
	claims["authorized"] = true
	claims["user"] = "NahtanN"

	tokenString, err := token.SignedString(byteSecret)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Error on login")
	}

	fmt.Fprintf(w, tokenString)
}
