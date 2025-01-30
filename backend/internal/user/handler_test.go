package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"web-chat/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Create(u user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockService) Authenticate(u user.User) (string, error) {
	args := m.Called(u)
	return args.String(0), args.Error(1)
}

func (m *MockService) GetUserByID(userID int) (user.User, error) {
	args := m.Called(userID)
	return args.Get(0).(user.User), args.Error(1)
}

func TestCreateUser_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockService)
	h := user.NewHandler(mockService)
	app.Post("/users", h.Create)

	requestBody := user.User{Name: "John Doe", Email: "john@example.com"}
	jsonBody, _ := json.Marshal(requestBody)

	mockService.On("Create", requestBody).Return(nil)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestCreateUser_InvalidRequest(t *testing.T) {
	app := fiber.New()
	mockService := new(MockService)
	h := user.NewHandler(mockService)
	app.Post("/users", h.Create)

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestLoginUser_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockService)
	h := user.NewHandler(mockService)
	app.Post("/login", h.Login)

	requestBody := user.User{Email: "john@example.com", Password: "secret"}
	jsonBody, _ := json.Marshal(requestBody)

	mockService.On("Authenticate", requestBody).Return("valid_token", nil)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestLoginUser_Unauthorized(t *testing.T) {
	app := fiber.New()
	mockService := new(MockService)
	h := user.NewHandler(mockService)
	app.Post("/login", h.Login)

	requestBody := user.User{Email: "john@example.com", Password: "wrong"}
	jsonBody, _ := json.Marshal(requestBody)

	mockService.On("Authenticate", requestBody).Return("", user.ErrUnauthorized)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	mockService.AssertExpectations(t)
}
