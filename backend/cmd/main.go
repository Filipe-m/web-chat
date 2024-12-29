package main

import (
	"fmt"
	"log"
	"os"
	"web-chat/internal/database"
	"web-chat/internal/middleware"
	"web-chat/internal/rooms"
	"web-chat/internal/user"

	"github.com/gofiber/contrib/swagger"

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

	roomRepo := rooms.NewRepository(db)
	roomService := rooms.NewService(roomRepo, userService)
	roomHandler := rooms.NewHandler(roomService)

	secret := os.Getenv("JWT")

	if secret == "" {
		fmt.Println("JWT secret não está definido no arquivo .env:", err)
		return
	}

	app := fiber.New()

	auth := middleware.JWTMiddleware(secret)
	authParam := middleware.JWTMiddlewareParam(secret)

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	cfg := swagger.Config{
		BasePath: "/",
		FilePath: "docs/swagger.json",
		Path:     "swagger",
		Title:    "Web Chat",
	}

	app.Use(swagger.New(cfg))
	
	app.Post("/auth/register", userHandler.Create)
	//TODO: verificar se retorna um 404 e 401
	app.Post("/auth/login", userHandler.Login)
	
	app.Post("/room", auth, roomHandler.Create)
	app.Get("/room", auth, roomHandler.GetRoom)
	app.Get("/room/:id", auth, roomHandler.GetRoom)
	//TODO: verificar se retorna um 404
	app.Delete("/room/:id", auth, roomHandler.Delete)

	app.Get("/ws/:id", authParam, websocket.New(roomHandler.Connect))
	//TODO: created_by retornando o id do usuário e não o nome
	app.Get("/messages/:id", auth, roomHandler.GetMessages)



	log.Fatal(app.Listen(":9090"))
}
