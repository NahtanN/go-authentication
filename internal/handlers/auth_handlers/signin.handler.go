package auth_handlers

import (
	"encoding/json"
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

type Tokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
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

	tokens, err := SignIn(handler.UserRepository, request)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, tokens)
}

func SignIn(userRepository database.UserRepository, request *SigninRequest) (*Tokens, error) {
	rows, err := userRepository.FindFirst(models.UserModel{
		Username: request.User,
		Email:    request.User,
	}).Select("id", "password").Exec()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id, password string

	for rows.Next() {
		err := rows.Scan(&id, &password)
		if err != nil {
			return nil, &utils.CustomError{
				Message: "Unable to parse data.",
			}
		}
	}

	defaultError := utils.CustomError{
		Message: "User or password invalid.",
	}

	if id == "" {
		return nil, &defaultError
	}

	match, err := utils.VerifyPassword(request.Password, password)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to validate user.",
		}
	}
	if !match {
		return nil, &defaultError
	}

	tokens, err := GenerateTokens(id)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to generate access token.",
		}
	}

	return tokens, nil
}

func GenerateTokens(id string) (*Tokens, error) {
	secret := os.Getenv("JWT_SECRET")
	byteSecret := []byte(secret)

	accessTokenClaims := jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(), // Expires in 10 minutes
		"id":  id,
	}
	accessTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	refreshTokenClaims := jwt.MapClaims{
		"exp": time.Now().Add(15 * time.Hour).Unix(),
	}
	refreshTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessToken, err := accessTokenString.SignedString(byteSecret)
	if err != nil {
		return nil, err
	}

	refreshToken, err := refreshTokenString.SignedString(byteSecret)
	if err != nil {
		return nil, err
	}

	tokens := Tokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return &tokens, nil
}
