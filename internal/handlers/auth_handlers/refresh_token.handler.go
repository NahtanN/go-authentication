package auth_handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nahtann/go-lab/internal/middlewares"
	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

type RefreshTokenHttpHandler struct {
	DB *pgxpool.Pool
}

type RefreshTokenRequest struct {
	Token string `json:"token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5..."`
}

// @Description	Creates a new pair of access and refresh tokens.
// @Tags			auth
// @Accept			json
// @Param			request	body	RefreshTokenRequest	true	"Request Body"
// @Produce		json
// @Success		201	{object}	auth_utils.Tokens
// @Failure		401	{object}	utils.CustomError	"Message: 'Invalid Request'"
// @router			/auth/refresh-token [post]
func NewRefreshTokenHttpHandler(db *pgxpool.Pool) *RefreshTokenHttpHandler {
	return &RefreshTokenHttpHandler{
		DB: db,
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

	tokens, err := RefreshToken(handler.DB, request.Token)
	if err != nil {
		return utils.WriteJSON(w, 401, err)
	}

	return utils.WriteJSON(w, http.StatusCreated, tokens)
}

func RefreshToken(
	db *pgxpool.Pool,
	tokenString string,
) (*auth_utils.Tokens, error) {
	token, valid := middlewares.ValidateJWT(tokenString)

	if !valid || !token.Valid {
		return nil, &utils.CustomError{
			Message: "Refresh Token not valid.",
		}
	}

	rows, err := db.Query(
		context.Background(),
		"SELECT id, user_id, used FROM refresh_tokens WHERE token = $1",
		tokenString,
	)
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
		_ = InvalidateUserRefreshTokens(db, userId)

		return nil, &utils.CustomError{
			Message: "Invalid Request",
		}
	}

	return UpdateUserRefreshToken(db, userId, id)
}

func InvalidateUserRefreshTokens(
	db *pgxpool.Pool,
	userId uint32,
) error {
	_, err := db.Exec(
		context.Background(),
		"UPDATE refresh_tokens SET used = true WHERE user_id = $1",
		userId,
	)
	if err != nil {
		return &utils.CustomError{
			Message: "Invalid Request",
		}
	}

	return nil
}

func UpdateUserRefreshToken(
	db *pgxpool.Pool,
	userId, parentTokenId uint32,
) (*auth_utils.Tokens, error) {
	tokens, err := auth_utils.CreateJwtTokens(userId)
	if err != nil {
		return nil, &utils.CustomError{
			Message: "Unable to generate access token.",
		}
	}

	_, err = db.Exec(
		context.Background(),
		"INSERT INTO refresh_tokens (parent_token_id, token, user_id, expires_at) VALUES ($1, $2, $3, $4)",
		parentTokenId,
		tokens.RefreshToken,
		userId,
		tokens.RefreshTokenExpiration,
	)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		context.Background(),
		"UPDATE refresh_tokens SET used = true WHERE id = $1",
		parentTokenId,
	)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
