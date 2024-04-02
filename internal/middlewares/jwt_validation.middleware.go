package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"

	"github.com/nahtann/go-authentication/internal/context_values"
	"github.com/nahtann/go-authentication/internal/utils"
)

type JWTValidationHttpHandler struct{}

func NewJWTValidationHttpHandler() *JWTValidationHttpHandler {
	return &JWTValidationHttpHandler{}
}

func (m *JWTValidationHttpHandler) Serve(
	w http.ResponseWriter,
	r *http.Request,
) (*http.Request, bool, error) {
	bearer := r.Header.Get("Authorization")
	token := strings.Split(bearer, " ")[1]

	tokenData, valid := ValidateJWT(token)

	if !valid || !tokenData.Valid {
		message := &utils.DefaultResponse{
			Message: "Not authorized.",
		}

		return nil, false, utils.WriteJSON(w, 401, message)
	}

	claims := tokenData.Claims.(jwt.MapClaims)
	userId := claims["id"]

	if userId == nil {
		message := &utils.DefaultResponse{
			Message: "Unable to parse claims.",
		}

		return nil, false, utils.WriteJSON(w, http.StatusInternalServerError, message)
	}

	ctx := context.WithValue(r.Context(), context_values.UserIdKey, userId)
	req := r.WithContext(ctx)

	return req, true, nil
}

func ValidateJWT(token string) (*jwt.Token, bool) {
	secret := os.Getenv("JWT_SECRET")

	tokenData, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("error on signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	return tokenData, err == nil
}
