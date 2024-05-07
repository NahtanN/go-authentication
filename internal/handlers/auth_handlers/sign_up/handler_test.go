package sign_up

import (
	"fmt"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) LIKE LOWER($1))").
		WithArgs("Test User").
		WillReturnRows(rows)

	rowsCopy := *rows

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) LIKE LOWER($1))").
		WithArgs("test@test.com").
		WillReturnRows(&rowsCopy)

	request := Request{
		Username: "Test User",
		Email:    "test@test.com",
		Password: "password",
	}

	mock.ExpectExec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)").
		WithArgs("Test User", "test@test.com", "asdf").
		WillReturnResult(pgxmock.NewResult("CREATE", 1))

	signUp := Handler{
		DB: mock,
		HashPassword: func(password string) (string, error) {
			return "asdf", nil
		},
	}

	result, err := signUp.Exec(&request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Message, "Sign up successfully")
}

func TestSignUpUsernameQuery(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) LIKE LOWER($1))").
		WithArgs("Test User").
		WillReturnError(fmt.Errorf("Some error"))

	request := Request{
		Username: "Test User",
		Email:    "test@test.com",
		Password: "password",
	}

	signUp := Handler{
		DB: mock,
	}

	result, err := signUp.Exec(&request)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Unable to validate user username.")
}

func TestSignUpEmailQuery(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) LIKE LOWER($1))").
		WithArgs("Test User").
		WillReturnRows(rows)

	rowsCopy := *rows

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) LIKE LOWER($1))").
		WithArgs("test@test.com").
		WillReturnRows(&rowsCopy)

	request := Request{
		Username: "Test User",
		Email:    "test@test.com",
		Password: "password",
	}

	mock.ExpectExec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)").
		WithArgs("Test User", "test@test.com", "asdf").
		WillReturnResult(pgxmock.NewResult("CREATE", 1))

	signUp := Handler{
		DB: mock,
	}

	result, err := signUp.Exec(&request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Message, "Sign up successfully")
}

func TestSignUp(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(username) LIKE LOWER($1))").
		WithArgs("Test User").
		WillReturnRows(rows)

	rowsCopy := *rows

	mock.ExpectQuery("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(email) LIKE LOWER($1))").
		WithArgs("test@test.com").
		WillReturnRows(&rowsCopy)

	request := Request{
		Username: "Test User",
		Email:    "test@test.com",
		Password: "password",
	}

	mock.ExpectExec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)").
		WithArgs("Test User", "test@test.com", "asdf").
		WillReturnResult(pgxmock.NewResult("CREATE", 1))

	signUp := Handler{
		DB: mock,
		HashPassword: func(password string) (string, error) {
			return "asdf", nil
		},
	}

	result, err := signUp.Exec(&request)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Message, "Sign up successfully")
}
