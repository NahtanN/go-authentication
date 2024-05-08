package wrappers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nahtann/go-lab/internal/utils"
)

type MockRequest struct {
	Data string `json:"data"`
}

type MockHandler struct{}

func (m *MockHandler) Exec(request *MockRequest) (*utils.DefaultResponse, error) {
	return &utils.DefaultResponse{
		Message: "Success.",
	}, nil
}

type MockHandlerFail struct{}

func (m *MockHandlerFail) Exec(request *MockRequest) (*string, error) {
	response := ""

	return &response, &utils.CustomError{
		Message: "Some Error.",
	}
}

func TestHttpWrapper(t *testing.T) {
	requestBody := `{"data":"example data"}`

	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	httpWrapper := HttpWrapper[MockRequest, utils.DefaultResponse]{
		Handler: &MockHandler{},
		ValidateRequest: func(s any) string {
			return ""
		},
	}

	err = httpWrapper.Serve(w, req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response utils.DefaultResponse

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, "Success.", response.Message)
}

func TestHttpWrapperDecodeError(t *testing.T) {
	requestBody := `"notValid":"example data"`

	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	httpWrapper := HttpWrapper[MockRequest, utils.DefaultResponse]{
		Handler: &MockHandler{},
	}

	err = httpWrapper.Serve(w, req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse utils.CustomError

	if err := json.NewDecoder(w.Body).Decode(&errorResponse); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, "Invalid Request.", errorResponse.Message)
}

func TestHttpWrapperValidateRequest(t *testing.T) {
	requestBody := `{"notValid":"example data"}`

	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	httpWrapper := HttpWrapper[MockRequest, utils.DefaultResponse]{
		Handler: &MockHandler{},
		ValidateRequest: func(s any) string {
			return "Some error message."
		},
	}

	err = httpWrapper.Serve(w, req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errorResponse utils.CustomError

	if err := json.NewDecoder(w.Body).Decode(&errorResponse); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, "Some error message.", errorResponse.Message)
}

func TestHttpWrapperExecError(t *testing.T) {
	requestBody := `{"data":"example data"}`

	req, err := http.NewRequest("POST", "/sign-in", bytes.NewBufferString(requestBody))
	if err != nil {
		t.Fatalf("an error '%s' was not expected while creating request", err)
	}
	w := httptest.NewRecorder()

	httpWrapper := HttpWrapper[MockRequest, string]{
		Handler: &MockHandlerFail{},
		ValidateRequest: func(s any) string {
			return ""
		},
	}

	err = httpWrapper.Serve(w, req)

	assert.Nil(t, err)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response utils.DefaultResponse

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	assert.Equal(t, "Some Error.", response.Message)
}
