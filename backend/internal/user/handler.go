package user

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service ServiceInterface
}

var (
	ErrUnauthorized       = errors.New("user unaunthorized")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("the given credentials are invalid")
)

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(c *fiber.Ctx) error {

	var request User
	if err := c.BodyParser(&request); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := h.service.Create(request)

	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (h *Handler) Login(c *fiber.Ctx) error {

	var request User
	if err := c.BodyParser(&request); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := h.service.Authenticate(request)

	if err != nil {
		switch {
		case errors.Is(err, ErrUnauthorized):
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
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
