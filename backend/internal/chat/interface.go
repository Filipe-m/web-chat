package chat

type RepositoryInterface interface {
	GetChatByID(chatID int, chat *Chat) error
	Create(chat *Chat, userId int) error
	Delete(chatID, userID int) error
}
