package chat

import (
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

var connections = struct {
	sync.Mutex
	clients map[*websocket.Conn]bool
}{
	clients: make(map[*websocket.Conn]bool),
}

func (h *Handler) Create(c *fiber.Ctx) error {

	var request Chat
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	userID, err := strconv.Atoi(c.Locals("id").(string))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	chat, err := h.service.Create(request, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(chat)
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	userID, err := strconv.Atoi(c.Locals("id").(string))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	chatIDStr := c.Params("id")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido, deve ser um número",
		})
	}

	err = h.service.Delete(chatID, userID)
	if err != nil {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) Connect(c *websocket.Conn) {
	var chat Chat

	userID := int(c.Locals("id").(float64))

	chatID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("chatID inválido:", err)
		c.Close()
	}

	chat, err = h.service.GetChatByID(chatID, chat)
	if err != nil {
		log.Println("Chat não encontrado:", err)
		c.Close()
	}

	h.service.Connect(c, userID, chat)
}
