package user

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{repository: repository}
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (s *Service) Authenticate(credentials User) string {

	user, err := s.repository.GetUserByEmail(credentials.Email)

	if (err != nil) || (!checkPasswordHash(credentials.Password, user.Password)) {
		fmt.Println(user)
		return ""
	}

	secret := os.Getenv("JWT")
	claims := jwt.MapClaims{
		"name": user.Name,
		"id":   user.ID,
		"exp":  time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return ""
	}

	return tokenString

}

func (s *Service) GetUserByID(userID int) (User, error) {
	var user User

	err := s.repository.GetUserByID(&user)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (s *Service) Create(user User) error {
	//validaçãow
	err := s.repository.Create(&user)

	if err != nil {
		return err
	}

	return nil
}
