package main

import (
	"log"
	"web-chat/cmd/handlers"
	"web-chat/internal/database"
	"web-chat/internal/user"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

func main() {

	db, err := database.Connect()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	userHandler := handlers.NewUserHandler(user.NewRepository(db))

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Post("/login", userHandler.Login)

	app.Get("/user/:id", userHandler.CreateUser)

	log.Fatal(app.Listen(":9090"))
}
