package rooms

import "web-chat/internal/user"

type RepositoryInterface interface {
	GetRoomByID(roomID int, room *Room) error
	Create(room *Room, userId int) error
	Delete(roomID, userID int) error
	SaveMessage(Message) (Message, error)
	GetRooms() ([]Room, error)
	GetAllMessages(roomId int) ([]Message, error)
	GetPaginatedMessages(roomId, page, size int) ([]Message, error)
}

type UserService interface {
	GetUserByID(id int) (user.User, error)
}
