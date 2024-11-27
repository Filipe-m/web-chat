package handlers

import (
	"web-chat/internal/user"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	repository *user.Repository
}

func NewUserHandler(repository *user.Repository) *User {
	return &User{repository: repository}
}

func (userHandler *User) CreateUser(c *fiber.Ctx) error {

	var request user.User
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	err := userHandler.repository.Create(&request)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (userHandler *User) Login(c *fiber.Ctx) error {

	var request user.User
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	token, err := userHandler.repository.Authenticate(request.Email, request.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Usuário ou senha inválidos"})
	}

	return c.Status(fiber.StatusOK).JSON(token)
}

func (userHandler *User) GetUserByID(userID int) user.User {
	var user user.User

	userHandler.repository.GetUserByID(&user)

	return user
}
