package user

import (
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *fiber.Ctx) error {

	var request User
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	err := h.service.Create(request)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (h *Handler) Login(c *fiber.Ctx) error {

	var request User
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	token := h.service.Authenticate(request)

	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Usuário ou senha inválidos"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (h *Handler) GetUserByID(userID int) (User, error) {

	user, err := h.service.GetUserByID(userID)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
