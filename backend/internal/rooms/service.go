package rooms

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type Service struct {
	repository  RepositoryInterface
	userService UserService
}

type userInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type broadCastMessage struct {
	Content string   `json:"content"`
	User    userInfo `json:"user"`
}
type storedMessages struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	User      userInfo  `json:"user"`
}

func NewService(repository RepositoryInterface, userService UserService) *Service {
	return &Service{
		repository:  repository,
		userService: userService,
	}
}

func (s *Service) Create(room Room, userID int) (Room, error) {

	err := s.repository.Create(&room, userID)

	if err != nil {
		return Room{}, err
	}

	return room, nil
}

func (s *Service) Delete(roomID int, userID int) error {

	err := s.repository.Delete(roomID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetRoomById(roomID int, room Room) (Room, error) {

	err := s.repository.GetRoomByID(roomID, &room)
	if err != nil {
		return Room{}, err
	}

	return room, nil
}

func (s *Service) SaveMessage(message Message) (Message, error) {
	message, err := s.repository.SaveMessage(message)

	if err != nil {
		return Message{}, err
	}
	return message, nil
}

func (s *Service) GetMessages(roomId, lastId, size int) ([]storedMessages, error) {
	var messages []storedMessages
	var err error

	if size == 0 {
		messages, err = s.repository.GetAllMessages(roomId)
		if err != nil {
			return []storedMessages{}, err
		}

		return messages, nil
	}
	messages, err = s.repository.GetPaginatedMessages(roomId, lastId, size)
	if err != nil {
		return []storedMessages{}, err
	}
	return messages, nil
}

func (s *Service) GetRooms() ([]Room, error) {

	rooms, err := s.repository.GetRooms()
	if err != nil {
		return []Room{}, err
	}
	return rooms, nil
}

func (s *Service) Connect(c *websocket.Conn, userID int, room Room) {

	user, err := s.userService.GetUserByID(userID)

	if err != nil {
		log.Println("Usuário inválido")
		return
	}

	log.Println("Conexão WebSocket estabelecida  \n\troom: ", room.Name, " \n\tusuário: ", user.Name)

	connections.Lock()
	connections.clients[c] = true
	connections.Unlock()

	defer func() {
		connections.Lock()
		delete(connections.clients, c)
		connections.Unlock()
		err := c.Close()
		if err != nil {
			log.Fatal("Não foi possível encerrar a conexão do websocket")
		}
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

		message := Message{
			Content:   string(msg),
			CreatedBy: userID,
			RoomId:    room.ID,
		}

		log.Println(message.Content)

		message, err = s.SaveMessage(message)

		if err != nil {
			log.Println(err)
		}

		userInfo := userInfo{
			Id:   user.ID,
			Name: user.Name,
		}

		broadcastMessage := broadCastMessage{
			Content: message.Content,
			User:    userInfo,
		}

		broadcast(broadcastMessage, c)
	}

}

func broadcast(messageObj broadCastMessage, sender *websocket.Conn) {
	connections.Lock()
	defer connections.Unlock()
	msgBytes, err := json.Marshal(messageObj)

	if err != nil {
		log.Println("Erro ao serializar o JSON")
	}

	for client := range connections.clients {
		if client != sender {
			err := client.WriteMessage(websocket.TextMessage, msgBytes)
			if err != nil {
				log.Printf("Erro ao enviar mensagem: %v", err)
				client.Close()
				delete(connections.clients, client)
			}
		}
	}
}
