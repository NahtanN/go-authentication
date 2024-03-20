package auth_handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"

	"github.com/nahtann/go-authentication/internal/storage/database"
	"github.com/nahtann/go-authentication/internal/storage/database/models"
	"github.com/nahtann/go-authentication/internal/utils"
)

type SignInHttpHandler struct {
	UserRepository database.UserRepository
}

type SigninRequest struct {
	User     string `json:"user"     validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SigninResponse struct {
	Token string `json:"token"`
}

func NewSignInHttpHandler(userRepository database.UserRepository) *SignInHttpHandler {
	return &SignInHttpHandler{
		UserRepository: userRepository,
	}
}

func (handler *SignInHttpHandler) Serve(w http.ResponseWriter, r *http.Request) error {
	request := new(SigninRequest)

	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		return utils.HttpServerInvalidRequest(w)
	}

	errorMessages := utils.Validate(request)
	if errorMessages != "" {
		message := utils.DefaultResponse{
			Message: errorMessages,
		}

		return utils.WriteJSON(w, http.StatusBadRequest, message)
	}

	token, err := SignIn(handler.UserRepository, request)
	if err != nil {
		return utils.WriteJSON(w, http.StatusInternalServerError, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, SigninResponse{
		Token: token,
	})
}

func SignIn(userRepository database.UserRepository, request *SigninRequest) (string, error) {
	result, err := userRepository.FindFirst(&models.UserModel{
		Username: request.User,
		Email:    request.User,
	}).Select(models.UserModel{}, "id", "password").Exec()
	if err != nil {
		return "", err
	}

	var id, password string

	err = result.Scan(&id, &password)
	if err != nil {
		fmt.Println(err)
		return "", &utils.CustomError{
			Message: "Unable to parse data.",
		}
	}

	fmt.Println(id, password)

	token, err := GenerateToken()
	if err != nil {
		return "", &utils.CustomError{
			Message: "Unable to generate access token",
		}
	}

	return token, nil
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
