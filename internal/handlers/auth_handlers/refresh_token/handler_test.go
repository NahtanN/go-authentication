package refresh_token

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

var mockToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

func mockValidToken() *jwt.Token {
	secret := []byte("secret")

	claims := jwt.MapClaims{
		"exp": time.Now().AddDate(0, 0, 15).Unix(),
	}
	tmp := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, _ := tmp.SignedString(secret)

	tokenData, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, fmt.Errorf("error on signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

	return tokenData
}

func TestRefreshToken(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	selectRows := mock.NewRows([]string{"id", "user_id", "used"}).
		AddRow(uint32(1), uint32(1), false)

	mock.ExpectQuery("SELECT id, user_id, used FROM refresh_tokens WHERE token = $1").
		WithArgs(mockToken).WillReturnRows(selectRows)

	handler := Handler{
		DB: mock,
		ValidateToken: func(token string) (*jwt.Token, bool) {
			return mockValidToken(), true
		},
		UpdateUserTokens: func(userId, parentTokenId uint32) (*auth_utils.Tokens, error) {
			return &auth_utils.Tokens{
				AccessToken:  "a",
				RefreshToken: "r",
			}, nil
		},
	}

	request := Request{
		Token: mockToken,
	}

	tokens, err := handler.Exec(&request)

	assert.Nil(t, err)
	assert.NotNil(t, tokens)
}

func TestRefreshTokenValidateTokenError(t *testing.T) {
	handler := Handler{
		ValidateToken: func(token string) (*jwt.Token, bool) {
			return mockValidToken(), false
		},
	}

	request := Request{
		Token: mockToken,
	}

	tokens, err := handler.Exec(&request)

	assert.Nil(t, tokens)
	assert.NotNil(t, err)
	assert.Equal(t, "Refresh Token not valid.", err.Error())
}

func TestRefreshTokenQueryError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	mock.ExpectQuery("SELECT id, user_id, used FROM refresh_tokens WHERE token = $1").
		WillReturnError(fmt.Errorf("Some error"))

	handler := Handler{
		DB: mock,
		ValidateToken: func(token string) (*jwt.Token, bool) {
			return mockValidToken(), true
		},
	}

	request := Request{
		Token: mockToken,
	}

	tokens, err := handler.Exec(&request)

	assert.Nil(t, tokens)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to validate refresh token data.", err.Error())
}

func TestRefreshTokenParseDataError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	selectRows := mock.NewRows([]string{"id", "user_id", "used"}).
		AddRow(1, 1, false)

	mock.ExpectQuery("SELECT id, user_id, used FROM refresh_tokens WHERE token = $1").
		WithArgs(mockToken).WillReturnRows(selectRows)

	handler := Handler{
		DB: mock,
		ValidateToken: func(token string) (*jwt.Token, bool) {
			return mockValidToken(), true
		},
	}

	request := Request{
		Token: mockToken,
	}

	tokens, err := handler.Exec(&request)

	assert.Nil(t, tokens)
	assert.NotNil(t, err)
	assert.Equal(t, "Unable to parse refresh token data.", err.Error())
}

func TestRefreshTokenInvalidateUserTokens(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	assert.Nil(t, err)
	defer mock.Close()

	selectRows := mock.NewRows([]string{"id", "user_id", "used"}).
		AddRow(uint32(1), uint32(1), true)

	mock.ExpectQuery("SELECT id, user_id, used FROM refresh_tokens WHERE token = $1").
		WithArgs(mockToken).WillReturnRows(selectRows)

	handler := Handler{
		DB: mock,
		ValidateToken: func(token string) (*jwt.Token, bool) {
			return mockValidToken(), true
		},
		InvalidateTokensByUser: func(userId uint32) error {
			return nil
		},
	}

	request := Request{
		Token: mockToken,
	}

	tokens, err := handler.Exec(&request)

	assert.Nil(t, tokens)
	assert.NotNil(t, err)
	assert.Equal(t, "Invalid Request.", err.Error())
}
