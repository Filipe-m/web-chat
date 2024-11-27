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

	userID, err := strconv.Atoi(c.Locals("id").(string))

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	err = chatHandler.repository.Create(&request, userID)
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

	if chats == nil {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.Status(fiber.StatusOK).JSON(chats)
}

func (chatHandler *Chat) Delete(c *fiber.Ctx) error {

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

	err = chatHandler.repository.Delete(chatID, userID)
	if err != nil {
		return c.SendStatus(fiber.StatusNoContent)
	}

	return c.SendStatus(fiber.StatusOK)
}
