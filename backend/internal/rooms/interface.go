package rooms

import (
	"web-chat/internal/user"

	"github.com/gofiber/contrib/websocket"
)

type RepositoryInterface interface {
	GetRoomByID(roomID int, room *Room) error
	Create(room *Room, userId int) error
	Delete(roomID, userID int) error
	SaveMessage(Message) (Message, error)
	GetRooms() ([]Room, error)
	GetAllMessages(roomId int) ([]storedMessages, error)
	GetPaginatedMessages(roomId, page, size int) ([]storedMessages, error)
}

type UserService interface {
	GetUserByID(id int) (user.User, error)
}

type ServiceInterface interface {
	Create(room Room, userID int) (Room, error)
	Delete(roomID int, userID int) error
	GetRoomById(roomID int, room Room) (Room, error)
	SaveMessage(message Message) (Message, error)
	GetMessages(roomId, lastId, size int) ([]storedMessages, error)
	GetRooms() ([]Room, error)
	Connect(c *websocket.Conn, userID int, room Room)
}
