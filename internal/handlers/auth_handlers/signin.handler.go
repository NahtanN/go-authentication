package auth_handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-lab/internal/utils"
)

type SignInHttpHandler struct {
	DB *pgxpool.Pool
}

type SigninRequest struct {
	User     string `json:"user"     validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Tokens struct {
	AccessToken            string    `json:"accessToken"`
	RefreshToken           string    `json:"refreshToken"`
	RefreshTokenExpiration time.Time `json:"-"`
}

func NewSignInHttpHandler(
	db *pgxpool.Pool,
) *SignInHttpHandler {
	return &SignInHttpHandler{
		DB: db,
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

	tokens, err := SignIn(handler.DB, request)
	if err != nil {
		return utils.WriteJSON(w, http.StatusBadRequest, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, tokens)
}

func SignIn(
	db *pgxpool.Pool,
	request *SigninRequest,
) (*Tokens, error) {
	rows, err := db.Query(
		context.Background(),
		"SELECT id, password FROM users WHERE email LIKE $1 OR username LIKE $1",
		request.User,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id uint32
	var password string

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

	if id == 0 {
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

	// Insert refresh token into database
	_, err = db.Exec(
		context.Background(),
		"INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)",
		tokens.RefreshToken,
		id,
		tokens.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func GenerateTokens(id uint32) (*Tokens, error) {
	secret := os.Getenv("JWT_SECRET")
	byteSecret := []byte(secret)

	accessTokenClaims := jwt.MapClaims{
		"exp": time.Now().Add(10 * time.Minute).Unix(), // Expires in 10 minutes
		"id":  id,
	}
	accessTokenString := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	refreshTokenExpiration := time.Now().AddDate(0, 0, 15)

	refreshTokenClaims := jwt.MapClaims{
		"exp": refreshTokenExpiration.Unix(),
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
		AccessToken:            accessToken,
		RefreshToken:           refreshToken,
		RefreshTokenExpiration: refreshTokenExpiration,
	}

	return &tokens, nil
}
