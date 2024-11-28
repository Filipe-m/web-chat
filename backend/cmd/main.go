package main

import (
	"fmt"
	"log"
	"os"
	"web-chat/internal/chat"
	"web-chat/internal/database"
	"web-chat/internal/middleware"
	"web-chat/internal/user"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

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

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)

	chatRepo := chat.NewRepository(db)
	chatService := chat.NewService(chatRepo, userService)
	chatHandler := chat.NewHandler(chatService)

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

	app.Get("/ws/:id", auth, websocket.New(chatHandler.Connect))

	app.Delete("/chat/:id", auth, chatHandler.Delete)

	app.Post("/chat", auth, chatHandler.Create)

	app.Post("/login", userHandler.Login)

	app.Post("/user", userHandler.Create)

	log.Fatal(app.Listen(":9090"))
}
