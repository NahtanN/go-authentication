package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/nahtann/go-authentication/internal/utils"
)

type SinginResponse struct {
	Token string `json:"token"`
}

func Signin(w http.ResponseWriter, r *http.Request) error {
	token, err := GenerateToken()
	if err != nil {
		apiError := utils.ApiError{
			Status:  http.StatusInternalServerError,
			Message: "Unable to sign in",
		}

		return utils.WriteJSON(w, http.StatusInternalServerError, apiError)
	}

	response := SinginResponse{
		Token: token,
	}

	return utils.WriteJSON(w, http.StatusCreated, response)
}

func GenerateToken() (string, error) {
	secret := os.Getenv("JWT_SECRET")
	byteSecret := []byte(secret)

	claims := jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(), // Expires in 10 minutes
		"id":  1,
	}
	tokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenString.SignedString(byteSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
