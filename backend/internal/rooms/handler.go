package rooms

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrNotFound = errors.New("content not found")
	//ErrUnauthorized = errors.New("not authorized")
	ErrNoContent = errors.New("no content")
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

	var request Room
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	userID := int(c.Locals("id").(float64))

	fmt.Println(userID)
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
	var room Room

	userID := int(c.Locals("id").(float64))

	roomID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("roomID inválido:", err)
		c.Close()
	}

	room, err = h.service.GetRoomById(roomID, room)
	if err != nil {
		log.Println("Chat não encontrado:", err)
		c.Close()
	}

	h.service.Connect(c, userID, room)
}

func (h *Handler) GetRoom(c *fiber.Ctx) error {

	var room Room

	roomIdParam := c.Params("id")

	if roomIdParam == "" {

		rooms, err := h.service.GetRooms()

		if err != nil {
			switch {
			case errors.Is(err, ErrNoContent):
				return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": err.Error()})
			default:
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
			}
		}

		return c.Status(fiber.StatusOK).JSON(rooms)
	}

	roomId, err := strconv.Atoi(roomIdParam)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	room, err = h.service.GetRoomById(roomId, room)

	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(room)

}

func (h *Handler) GetMessages(c *fiber.Ctx) error {
	var page, size int

	roomIdParam := c.Params("id")
	pageParam := c.Query("page")
	sizeParam := c.Query("size")

	if roomIdParam == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	roomId, _ := strconv.Atoi(roomIdParam)

	if pageParam == "" || sizeParam == "" {
		page = 0
		size = 0
	} else {
		page, _ = strconv.Atoi(pageParam)
		size, _ = strconv.Atoi(sizeParam)
		if page <= 0 || size <= 0 {
			page = 1
			size = 10
		}
	}

	messages, err := h.service.GetMessages(roomId, page, size)

	if err != nil {
		switch {
		case errors.Is(err, ErrNoContent):
			return c.Status(fiber.StatusNoContent).JSON(fiber.Map{"message": err.Error()})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
		}
	}

	return c.Status(fiber.StatusOK).JSON(messages)
}
