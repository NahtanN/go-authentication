package sign_in

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"

	"github.com/nahtann/go-lab/internal/utils"
	auth_utils "github.com/nahtann/go-lab/internal/utils/auth"
)

type MockHandler struct{}

func (m *MockHandler) Exec(request *SigninRequest) (*auth_utils.Tokens, error) {
	return &auth_utils.Tokens{
		AccessToken:  "a",
		RefreshToken: "r",
	}, nil
}

type MockHandlerFail struct{}

func (m *MockHandlerFail) Exec(request *SigninRequest) (*auth_utils.Tokens, error) {
	return nil, &utils.CustomError{
		Message: "Some Error.",
	}
}

func TestHttpWrapper(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	requestBody := `{"user":"example@example.com","password":"password123"}`

	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	signInHttpWrapper := HttpWrapper{
		Handler: &MockHandler{},
		ValidateRequest: func(s any) string {
			return ""
		},
	}

	err = signInHttpWrapper.Serve(w, req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)

	var tokens auth_utils.Tokens

	if err := json.NewDecoder(w.Body).Decode(&tokens); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	assert.Equal(
		t, "a", tokens.AccessToken,
	)

	assert.Equal(
		t, "r", tokens.RefreshToken,
	)
}

func TestHttpWrapperError(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	requestBody := `{"a":"example@example.com","b":"password123"}`

	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	signInHttpWrapper := HttpWrapper{
		Handler: &MockHandlerFail{},
		ValidateRequest: func(s any) string {
			return ""
		},
	}

	err = signInHttpWrapper.Serve(w, req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	assert.Equal(
		t, "Some Error.", body.Message,
	)
}

func TestHttpWrapperInvalidRequest(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(""))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	signInHttpWrapper := HttpWrapper{}

	err = signInHttpWrapper.Serve(w, req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, "Invalid Request.", body.Message)
}

func TestHttpWrapperFieldValidation(t *testing.T) {
	mock, err := pgxmock.NewPool(pgxmock.QueryMatcherOption(pgxmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mock.Close()

	requestBody := `{"a":"example@example.com","b":"password123"}`

	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	signInHttpWrapper := HttpWrapper{
		ValidateRequest: func(s any) string {
			return "Field `User` failed validation. Field `Password` failed validation."
		},
	}

	err = signInHttpWrapper.Serve(w, req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var body struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	assert.Equal(
		t,
		"Field `User` failed validation. Field `Password` failed validation.",
		body.Message,
	)
}
