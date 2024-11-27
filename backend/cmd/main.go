package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"web-chat/cmd/handlers"
	"web-chat/internal/chat"
	"web-chat/internal/database"
	"web-chat/internal/middleware"
	"web-chat/internal/user"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var connections = struct {
	sync.Mutex
	clients map[*websocket.Conn]bool
}{
	clients: make(map[*websocket.Conn]bool),
}

func broadcast(msg []byte, sender *websocket.Conn) {
	connections.Lock()
	defer connections.Unlock()

	for client := range connections.clients {
		if client != sender { // Evita ecoar para o cliente que enviou
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Erro ao enviar mensagem: %v", err)
				client.Close()
				delete(connections.clients, client) // Remove conexões quebradas
			}
		}
	}
}

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		fmt.Printf("erro ao carregar o arquivo .env: %v", err)
		return
	}

	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	userHandler := handlers.NewUserHandler(user.NewRepository(db))
	chatHandler := handlers.NewChatHandler(chat.NewRepository(db))

	secret := os.Getenv("JWT")

	if secret == "" {
		fmt.Println("JWT secret não está definido no arquivo .env:", err)
		return
	}

	app := fiber.New()

	auth := middleware.JWTMiddleware(secret)

	app.Get("/", auth, func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// Rota WebSocket para comunicação
	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		roomId := c.Params("id")
		fmt.Println(roomId)
		// Adiciona a conexão ao mapa
		connections.Lock()
		connections.clients[c] = true
		connections.Unlock()

		fmt.Println("Nova conexão estabelecida")

		defer func() {
			// Remove a conexão ao desconectar
			connections.Lock()
			delete(connections.clients, c)
			connections.Unlock()
			fmt.Println("Conexão encerrada")
			c.Close()
		}()

		// Escuta mensagens do cliente
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Printf("Erro ao ler mensagem: %v", err)
				break
			}
			log.Printf("Mensagem recebida: %s", msg)

			// Envia a mensagem para todos os outros clientes conectados
			broadcast(msg, c)
		}
	}))

	app.Get("/chat", auth, chatHandler.GetChats)

	app.Delete("/chat/:id", auth, chatHandler.Delete)

	app.Post("/chat", auth, chatHandler.Create)

	app.Post("/login", userHandler.Login)

	app.Post("/user", userHandler.CreateUser)

	log.Fatal(app.Listen(":9090"))
}
