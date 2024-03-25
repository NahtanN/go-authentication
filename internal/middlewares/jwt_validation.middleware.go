package middlewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/nahtann/go-authentication/internal/utils"
)

type JWTValidationHttpHandler struct{}

func NewJWTValidationHttpHandler() *JWTValidationHttpHandler {
	return &JWTValidationHttpHandler{}
}

func (m *JWTValidationHttpHandler) Serve(w http.ResponseWriter, r *http.Request) (bool, error) {
	bearer := r.Header.Get("Authorization")
	token := strings.Split(bearer, " ")[1]

	valid := ValidateJWT(token)

	if !valid {
		message := &utils.DefaultResponse{
			Message: "Not authorized.",
		}

		return false, utils.WriteJSON(w, 401, message)
	}

	return true, nil
}

func ValidateJWT(token string) bool {
	secret := os.Getenv("JWT_SECRET")

	_, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("Error on signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	return err == nil
}
