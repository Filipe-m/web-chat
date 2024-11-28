package chat

import "web-chat/internal/user"

type RepositoryInterface interface {
	GetChatByID(chatID int, chat *Chat) error
	Create(chat *Chat, userId int) error
	Delete(chatID, userID int) error
}

type UserService interface {
	GetUserByID(id int) (user.User, error)
}
