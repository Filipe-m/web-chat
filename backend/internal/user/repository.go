package user

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (repo *Repository) Authenticate(email, password string) (string, error) {
	query := "SELECT id, name, password FROM users WHERE email = $1"
	row := repo.DB.QueryRow(query, email)

	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	if !checkPasswordHash(password, user.Password) {
		return "", nil
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
		return "", err
	}

	return tokenString, nil
}

func (repo *Repository) Create(user *User) error {

	password, err := hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("erro ao criptografar a senha: %v", err)
	}

	query := "INSERT INTO users (name, password, email) VALUES ($1, $2, $3)"

	_, err = repo.DB.Exec(query, user.Name, password, user.Email)

	if err != nil {
		return fmt.Errorf("erro ao inserir usuário: %v", err)
	}

	return nil
}

func (repo *Repository) GetUserByID(user *User) error {
	query := "SELECT id, name, password, email, created_at, updated_at FROM users WHERE id = $1 RETURNING"

	err := repo.DB.QueryRow(query, user.ID).Scan(&user.ID, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("usuário não encontrado: %v", err)
	}

	return nil

}
