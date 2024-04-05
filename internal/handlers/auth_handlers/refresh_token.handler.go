package auth_handlers

import (
	"encoding/json"
	"net/http"

	"github.com/nahtann/go-authentication/internal/middlewares"
	"github.com/nahtann/go-authentication/internal/storage/database"
	"github.com/nahtann/go-authentication/internal/storage/database/models"
	"github.com/nahtann/go-authentication/internal/storage/database/query_builder"
	"github.com/nahtann/go-authentication/internal/utils"
)

type RefreshTokenHttpHandler struct {
	RefreshTokenRepository database.RefreshTokenRepository
}

type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required"`
}

func NewRefreshTokenHttpHandler(
	refreshTokenRepository database.RefreshTokenRepository,
) *RefreshTokenHttpHandler {
	return &RefreshTokenHttpHandler{
		RefreshTokenRepository: refreshTokenRepository,
	}
}

func (handler *RefreshTokenHttpHandler) Serve(w http.ResponseWriter, r *http.Request) error {
	request := new(RefreshTokenRequest)

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

	tokens, err := RefreshToken(handler.RefreshTokenRepository, request.Token)
	if err != nil {
		return utils.WriteJSON(w, 401, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, tokens)
}

func RefreshToken(
	refreshTokenRepository database.RefreshTokenRepository,
	tokenString string,
) (*Tokens, error) {
	token, valid := middlewares.ValidateJWT(tokenString)

	if !valid || !token.Valid {
		return nil, &utils.CustomError{
			Message: "Refresh Token not valid.",
		}
	}

	rows, err := refreshTokenRepository.FindFirst().
		Where(
			query_builder.Equals("token", tokenString),
		).
		Select("id", "user_id", "used").
		Exec()
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to validate refresh token data.",
		}
	}
	defer rows.Close()

	var id, userId uint32
	var used bool

	for rows.Next() {
		err := rows.Scan(&id, &userId, &used)
		if err != nil {
			return nil, &utils.CustomError{
				Message: "Unable to parse refresh token data.",
			}
		}
	}

	if used || userId == 0 {
		InvalidateUserRefreshTokens(userId)

		return nil, &utils.CustomError{
			Message: "Invalid Request",
		}
	}

	return UpdateUserRefreshToken(refreshTokenRepository, userId, id)
}

func InvalidateUserRefreshTokens(userId uint32) {
}

func UpdateUserRefreshToken(
	refreshTokenRepository database.RefreshTokenRepository,
	userId, parentTokenId uint32,
) (*Tokens, error) {
	tokens, err := GenerateTokens(userId)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to generate access token.",
		}
	}

	err = refreshTokenRepository.Create(models.RefreshTokenModel{
		ParentTokenId: parentTokenId,
		Token:         tokens.RefreshToken,
		UserId:        userId,
		ExpiresAt:     tokens.RefreshTokenExpiration,
	})
	if err != nil {
		return nil, err
	}

	err = refreshTokenRepository.Update(models.RefreshTokenModel{
		Id:   parentTokenId,
		Used: true,
	})
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
