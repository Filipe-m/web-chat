package handlers

import (
	"log"
	"strconv"
	"sync"
	"web-chat/internal/chat"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type Chat struct {
	repository *chat.Repository
}

var connections = struct {
	sync.Mutex
	clients map[*websocket.Conn]bool
}{
	clients: make(map[*websocket.Conn]bool),
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

func (chatHandler *Chat) Connect(c *websocket.Conn) {
	var chat chat.Chat

	chatID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Println("ID inválido:", err)
		c.Close()
		return
	}

	err = chatHandler.repository.GetChatByID(chatID, &chat)
	if err != nil {
		log.Println("Chat não encontrado:", err)
		c.Close()
		return
	}

	log.Println("Conexão WebSocket estabelecida para o chat ", chat.Name, "")

	connections.Lock()
	connections.clients[c] = true
	connections.Unlock()

	defer func() {
		connections.Lock()
		delete(connections.clients, c)
		connections.Unlock()
		c.Close()
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Erro ao ler mensagem: %v", err)
			break
		}

		broadcast(msg, c)
	}
}

func broadcast(msg []byte, sender *websocket.Conn) {
	connections.Lock()
	defer connections.Unlock()

	for client := range connections.clients {
		if client != sender {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Erro ao enviar mensagem: %v", err)
				client.Close()
				delete(connections.clients, client)
			}
		}
	}
}
