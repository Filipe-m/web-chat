package user_test

import (
	"errors"
	"os"
	"testing"
	"web-chat/internal/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByEmail(email string) (user.User, error) {
	args := m.Called(email)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *MockRepository) GetUserByID(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func (m *MockRepository) Create(u *user.User) error {
	args := m.Called(u)
	return args.Error(0)
}

func hashPassword(password string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed)
}

func TestAuthenticate_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := user.NewService(mockRepo)

	hashedPassword := hashPassword("secret")
	expectedUser := user.User{ID: 1, Name: "John Doe", Email: "john@example.com", Password: hashedPassword}

	mockRepo.On("GetUserByEmail", "john@example.com").Return(expectedUser, nil)

	os.Setenv("JWT", "test_secret")

	token, err := service.Authenticate(user.User{Email: "john@example.com", Password: "secret"})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedToken, _ := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test_secret"), nil
	})

	assert.True(t, parsedToken.Valid)
	mockRepo.AssertExpectations(t)
}

func TestAuthenticate_InvalidPassword(t *testing.T) {
	mockRepo := new(MockRepository)
	service := user.NewService(mockRepo)

	expectedUser := user.User{ID: 1, Name: "John Doe", Email: "john@example.com", Password: hashPassword("secret")}

	mockRepo.On("GetUserByEmail", "john@example.com").Return(expectedUser, nil)

	token, err := service.Authenticate(user.User{Email: "john@example.com", Password: "wrong"})

	assert.ErrorIs(t, err, user.ErrUnauthorized)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_Success_Service(t *testing.T) {
	mockRepo := new(MockRepository)
	service := user.NewService(mockRepo)

	newUser := user.User{Name: "John Doe", Email: "john@example.com", Password: "password"}

	mockRepo.On("Create", mock.AnythingOfType("*user.User")).Return(nil)

	err := service.Create(newUser)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_InvalidEmail(t *testing.T) {
	mockRepo := new(MockRepository)
	service := user.NewService(mockRepo)

	newUser := user.User{Name: "John Doe", Email: "invalid-email", Password: "password"}

	err := service.Create(newUser)

	assert.ErrorIs(t, err, user.ErrInvalidCredentials)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetUserByID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := user.NewService(mockRepo)

	expectedUser := user.User{ID: 1, Name: "John Doe", Email: "john@example.com"}
	mockRepo.On("GetUserByID", mock.AnythingOfType("*user.User")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*user.User)
		*arg = expectedUser
	})

	userData, err := service.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, userData)
	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := user.NewService(mockRepo)

	mockRepo.On("GetUserByID", mock.AnythingOfType("*user.User")).Return(errors.New("user not found"))

	userData, err := service.GetUserByID(1)

	assert.Error(t, err)
	assert.Equal(t, user.User{}, userData)
	mockRepo.AssertExpectations(t)
}
