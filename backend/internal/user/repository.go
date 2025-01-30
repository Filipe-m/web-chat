package user

import (
	"database/sql"
	"fmt"

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

func (repo *Repository) Create(user *User) error {

	password, err := hashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("erro ao criptografar a senha: %v", err)
	}

	query := "INSERT INTO users (name, password, email) VALUES ($1, $2, $3)"

	_, err = repo.DB.Exec(query, user.Name, password, user.Email)

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetUserByID(user *User) error {
	query := "SELECT id, name, password, email, created_at, updated_at FROM users WHERE id = $1"

	err := repo.DB.QueryRow(query, user.ID).Scan(&user.ID, &user.Name, &user.Password, &user.Email, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return ErrUserNotFound
		}
		return err
	}

	return nil

}

func (repo *Repository) GetUserByEmail(email string) (User, error) {
	var user User

	query := "SELECT id, name, email, password, created_at, updated_at FROM users WHERE email = $1"

	err := repo.DB.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return User{}, fmt.Errorf("usuário não encontrado: %v", err)
		}
		return User{}, err
	}

	return user, nil
}
