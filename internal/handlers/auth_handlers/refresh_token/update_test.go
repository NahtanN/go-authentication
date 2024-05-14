package refresh_token

import (
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

var (
	parentTokenId = uint32(1)
	userId        = uint32(1)
)

var tokens = auth_utils.Tokens{
	AccessToken:            "a",
	RefreshToken:           "r",
	RefreshTokenExpiration: time.Now(),
}

func TestUpdateHandler(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO refresh_tokens (parent_token_id, token, user_id, expires_at) VALUES ($1, $2, $3, $4)").
		WithArgs(parentTokenId, tokens.RefreshToken, userId, tokens.RefreshTokenExpiration).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectExec("UPDATE refresh_tokens SET used = true WHERE id = $1").
		WithArgs(parentTokenId).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectCommit()

	handler := UpdateHandler{
		DB: mock,
		CreateTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return &tokens, nil
		},
	}

	result, err := handler.UserTokens(userId, parentTokenId)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.AccessToken, tokens.AccessToken)
	assert.Equal(t, result.RefreshToken, tokens.RefreshToken)
}

func TestUpdateHandlerCreateTokensError(t *testing.T) {
	handler := UpdateHandler{
		CreateTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return nil, fmt.Errorf("Some error")
		},
	}

	result, err := handler.UserTokens(userId, parentTokenId)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to create tokens.", err.Error())
}

func TestUpdateHandlerStartTransactionError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	mock.ExpectBegin().WillReturnError(fmt.Errorf("Some error"))

	handler := UpdateHandler{
		DB: mock,
		CreateTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return &tokens, nil
		},
	}

	result, err := handler.UserTokens(userId, parentTokenId)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to start transaction.", err.Error())
}

func TestUpdateHandlerSaveTokensError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO refresh_tokens (parent_token_id, token, user_id, expires_at) VALUES ($1, $2, $3, $4)").
		WithArgs(parentTokenId, tokens.RefreshToken, userId, tokens.RefreshTokenExpiration).
		WillReturnError(fmt.Errorf("Some error"))

	mock.ExpectRollback()

	handler := UpdateHandler{
		DB: mock,
		CreateTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return &tokens, nil
		},
	}

	result, err := handler.UserTokens(userId, parentTokenId)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to save tokens.", err.Error())
}

func TestUpdateHandlerUpdateRefreshTokenError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO refresh_tokens (parent_token_id, token, user_id, expires_at) VALUES ($1, $2, $3, $4)").
		WithArgs(parentTokenId, tokens.RefreshToken, userId, tokens.RefreshTokenExpiration).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectExec("UPDATE refresh_tokens SET used = true WHERE id = $1").
		WithArgs(parentTokenId).
		WillReturnError(fmt.Errorf("Some error"))

	mock.ExpectRollback()

	handler := UpdateHandler{
		DB: mock,
		CreateTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return &tokens, nil
		},
	}

	result, err := handler.UserTokens(userId, parentTokenId)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to update refresh token.", err.Error())
}

func TestUpdateHandlerCommitTransactionError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO refresh_tokens (parent_token_id, token, user_id, expires_at) VALUES ($1, $2, $3, $4)").
		WithArgs(parentTokenId, tokens.RefreshToken, userId, tokens.RefreshTokenExpiration).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectExec("UPDATE refresh_tokens SET used = true WHERE id = $1").
		WithArgs(parentTokenId).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectCommit().WillReturnError(fmt.Errorf("Some error"))

	handler := UpdateHandler{
		DB: mock,
		CreateTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return &tokens, nil
		},
	}

	result, err := handler.UserTokens(userId, parentTokenId)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to commit transaction.", err.Error())
}
