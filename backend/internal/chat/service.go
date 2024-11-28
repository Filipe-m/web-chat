package chat

import (
	"log"

	"github.com/gofiber/contrib/websocket"
)

type Service struct {
	repository  RepositoryInterface
	userService UserService
}

func NewService(repository RepositoryInterface, userService UserService) *Service {
	return &Service{
		repository:  repository,
		userService: userService,
	}
}

func (s *Service) Create(chat Chat, userID int) (Chat, error) {

	err := s.repository.Create(&chat, userID)

	if err != nil {
		return Chat{}, err
	}

	return chat, nil
}

func (s *Service) Delete(chatID int, userID int) error {

	err := s.repository.Delete(chatID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetChatByID(chatID int, chat Chat) (Chat, error) {

	err := s.repository.GetChatByID(chatID, &chat)
	if err != nil {
		return Chat{}, err
	}

	return chat, nil
}

func (s *Service) Connect(c *websocket.Conn, userID int, chat Chat) {

	user, err := s.userService.GetUserByID(userID)

	if err != nil {
		log.Println("Usuário inválido")
		return
	}

	log.Println("Conexão WebSocket estabelecida  \n\tchat: ", chat.Name, " \n\tusuário: ", user.Name)

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

			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseNoStatusReceived) {
				log.Println("Conexão fechada pelo cliente.")
			} else {
				log.Printf("Erro ao ler mensagem: %v", err)
			}
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
