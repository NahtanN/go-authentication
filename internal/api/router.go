package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

type Message struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

type EndPoint func(w http.ResponseWriter, r *http.Request) error

func SetRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /public", public)
	mux.HandleFunc("POST /private", handleAuth(private))
	mux.HandleFunc("POST /login", login)
}

func public(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is a public route")
}

func private(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var message Message

	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		return err
	}

	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		return err
	}

	return nil
}

func login(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("JWT_SECRET")
	byteSecret := []byte(secret)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["user"] = "NahtanN"

	tokenString, err := token.SignedString(byteSecret)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Error on login")
	}

	fmt.Fprintf(w, tokenString)
}

func handleAuth(endpoint EndPoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("JWT_SECRET")
		authorization := r.Header.Get("Authorization")

		if len(authorization) == 0 {
			fmt.Fprintf(w, "token not provided")
			return
		}

		tokenString := strings.Split(authorization, " ")[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				w.WriteHeader(http.StatusUnauthorized)

				_, err := w.Write([]byte("You`re Unauthorized"))
				if err != nil {
					return nil, err
				}
			}

			return []byte(secret), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			w.Write([]byte("You're Unauthorized due to error parsing the JWT"))
			return
		}

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)

			fmt.Println(claims["user"])

			if err := endpoint(w, r); err != nil {
				fmt.Fprintf(w, "Request Error")
				return
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)

			_, err := w.Write([]byte("You're Unauthorized due to invalid token"))
			if err != nil {
				return
			}
		}
	}
}
