package handlers

import (
	"strconv"
	"web-chat/internal/chat"

	"github.com/gofiber/fiber/v2"
)

type Chat struct {
	repository *chat.Repository
}

func NewChatHandler(repository *chat.Repository) *Chat {
	return &Chat{repository: repository}
}

func (chatHandler *Chat) Create(c *fiber.Ctx) error {

	var request chat.Chat
	if err := c.BodyParser(&request); err != nil {
		return err
	}

	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("userID inválido")
	}

	err := chatHandler.repository.Create(&request, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(request)
}

func (chatHandler *Chat) GetChats(c *fiber.Ctx) error {

	chats, err := chatHandler.repository.GetChats()
	if err != nil {
		return nil
	}

	return c.Status(fiber.StatusOK).JSON(chats)
}

func (chatHandler *Chat) Delete(c *fiber.Ctx) error {

	userID, ok := c.Locals("userID").(int)
	if !ok {
		return c.Status(fiber.StatusBadRequest).SendString("userID inválido")
	}

	chatIDStr := c.Params("chatID")
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID inválido, deve ser um número",
		})
	}

	err = chatHandler.repository.Delete(chatID, userID)
	if err != nil {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.SendStatus(fiber.StatusOK)
}
