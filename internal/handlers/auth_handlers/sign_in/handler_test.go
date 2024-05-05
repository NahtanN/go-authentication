package sign_in

import (
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

var request = SigninRequest{
	User:     "nahtann@outlook.com",
	Password: "password",
}

var tokens = auth_utils.Tokens{
	AccessToken:            "asdf",
	RefreshToken:           "asfd",
	RefreshTokenExpiration: time.Now(),
}

func TestSignInShouldReturnTokens(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	var id uint32 = 1

	rows := mock.NewRows([]string{"id", "password"}).
		AddRow(id, "password")

	mock.ExpectQuery("SELECT id, password FROM users WHERE email LIKE $1 OR username LIKE $1").
		WithArgs(request.User).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)").
		WithArgs(tokens.RefreshToken, id, tokens.RefreshTokenExpiration).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	signIn := Handler{
		DB: mock,
		VerifyPassword: func(password, hashedPassword string) (bool, error) {
			return true, nil
		},
		CreateJwtTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return &tokens, nil
		},
	}

	signInTokens, err := signIn.Exec(&request)

	assert.NotNil(t, signInTokens)
	assert.Nil(t, err)
	assert.EqualValues(t, &tokens, signInTokens)
}

func TestSignInInvalidPassword(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	var id uint32 = 1

	rows := mock.NewRows([]string{"id", "password"}).
		AddRow(id, "password")

	mock.ExpectQuery("SELECT id, password FROM users WHERE email LIKE $1 OR username LIKE $1").
		WithArgs(request.User).
		WillReturnRows(rows)

	signIn := Handler{
		DB: mock,
		VerifyPassword: func(password, hashedPassword string) (bool, error) {
			return false, nil
		},
	}

	signInTokens, err := signIn.Exec(&request)

	assert.Nil(t, signInTokens)
	assert.NotNil(t, err)
	assert.Equal(t, "User or password invalid.", err.Error())
}

func TestSignInCreateTokensError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	var id uint32 = 1

	rows := mock.NewRows([]string{"id", "password"}).
		AddRow(id, "password")

	mock.ExpectQuery("SELECT id, password FROM users WHERE email LIKE $1 OR username LIKE $1").
		WithArgs(request.User).
		WillReturnRows(rows)

	signIn := Handler{
		DB: mock,
		VerifyPassword: func(password, hashedPassword string) (bool, error) {
			return true, nil
		},
		CreateJwtTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return nil, fmt.Errorf("Some error")
		},
	}

	signInTokens, err := signIn.Exec(&request)

	assert.NotNil(t, err)
	assert.Nil(t, signInTokens)
	assert.Equal(t, "Unable to create access token.", err.Error())
}

func TestSignInShouldOnUpdateRefreshToken(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	var id uint32 = 1

	rows := mock.NewRows([]string{"id", "password"}).
		AddRow(id, "password")

	mock.ExpectQuery("SELECT id, password FROM users WHERE email LIKE $1 OR username LIKE $1").
		WithArgs(request.User).
		WillReturnRows(rows)

	mock.ExpectExec("INSERT INTO refresh_tokens (token, user_id, expires_at) VALUES ($1, $2, $3)").
		WithArgs(tokens.RefreshToken, id, tokens.RefreshTokenExpiration).
		WillReturnError(fmt.Errorf("Some error"))

	signIn := Handler{
		DB: mock,
		VerifyPassword: func(password, hashedPassword string) (bool, error) {
			return true, nil
		},
		CreateJwtTokens: func(id uint32) (*auth_utils.Tokens, error) {
			return &tokens, nil
		},
	}

	signInTokens, err := signIn.Exec(&request)

	assert.NotNil(t, err)
	assert.Nil(t, signInTokens)
	assert.Equal(t, "Internal error.", err.Error())
}
