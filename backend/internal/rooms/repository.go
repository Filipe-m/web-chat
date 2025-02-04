package rooms

import (
	"database/sql"
	"errors"
	"fmt"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (repo *Repository) Create(room *Room, userID int) error {

	query := "INSERT INTO rooms (name, created_by) VALUES ($1, $2) RETURNING id, name, created_by, created_at, updated_at"

	err := repo.DB.QueryRow(query, room.Name, userID).Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		return fmt.Errorf("erro ao criar rooms: %v", err)
	}

	return nil
}

func (repo *Repository) GetRooms() ([]Room, error) {

	query := "SELECT id, name, created_by, created_at, updated_at FROM rooms"

	var rooms []Room

	rows, err := repo.DB.Query(query)

	if err != nil {
		return nil, fmt.Errorf("erro ao executar a query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var room Room

		err := rows.Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("erro ao escanear os resultados: %w", err)
		}
		rooms = append(rooms, room)
	}

	if len(rooms) == 0 {
		return nil, ErrNoContent
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar os resultados: %w", err)
	}

	return rooms, nil

}

func (repo *Repository) Delete(roomID, userID int) error {
	query := "DELETE FROM rooms WHERE created_by = $1 AND id = $2"

	result, err := repo.DB.Exec(query, userID, roomID)
	if err != nil {
		return fmt.Errorf("erro ao executar a query de deleção: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return ErrUnauthorized
	}

	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (repo *Repository) GetRoomByID(roomID int, room *Room) error {

	query := "SELECT id, name, created_by, created_at, updated_at FROM rooms where id = $1"

	err := repo.DB.QueryRow(query, roomID).Scan(&room.ID, &room.Name, &room.CreatedBy, &room.CreatedAt, &room.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
	}

	return nil
}

func (repo *Repository) SaveMessage(message Message) (Message, error) {
	query := "INSERT INTO messages (content, created_by, room_id) VALUES ($1, $2, $3) RETURNING id, content, created_by, room_id, created_at, updated_at"

	err := repo.DB.QueryRow(query, message.Content, message.CreatedBy, message.RoomId).Scan(&message.ID, &message.Content, &message.CreatedBy, &message.RoomId, &message.CreatedAt, &message.UpdatedAt)

	if err != nil {
		return Message{}, fmt.Errorf("erro ao criar rooms: %v", err)
	}

	return message, nil
}

func (repo *Repository) GetAllMessages(roomID int) ([]storedMessages, error) {

	query := "SELECT m.id, m.content , m.created_at, m.updated_at, u.id, u.name FROM messages as m join users as u on u.id = m.created_by where room_id = $1 order by created_at desc"

	var messages []storedMessages

	rows, err := repo.DB.Query(query, roomID)

	if err != nil {
		return nil, fmt.Errorf("erro ao executar a query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var message storedMessages

		err := rows.Scan(&message.Id, &message.Content, &message.CreatedAt, &message.UpdatedAt, &message.User.Id, &message.User.Name)

		if err != nil {
			return nil, fmt.Errorf("erro ao escanear os resultados: %w", err)
		}

		messages = append(messages, message)
	}

	if len(messages) == 0 {
		return nil, ErrNoContent
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar os resultados: %w", err)
	}

	return messages, nil

}

func (repo *Repository) GetPaginatedMessages(roomId, lastId, size int) ([]storedMessages, error) {

	query := "SELECT m.id, m.content , m.created_at, m.updated_at, u.id, u.name FROM messages as m join users as u on u.id = m.created_by WHERE m.room_id = $1 AND m.ID > $2 ORDER BY m.ID DESC LIMIT $3"

	var messages []storedMessages

	rows, err := repo.DB.Query(query, roomId, lastId, size)

	if err != nil {
		return nil, fmt.Errorf("erro ao executar a query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var message storedMessages

		err := rows.Scan(&message.Id, &message.Content, &message.CreatedAt, &message.UpdatedAt, &message.User.Id, &message.User.Name)

		if err != nil {
			return nil, fmt.Errorf("erro ao escanear os resultados: %w", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar os resultados: %w", err)
	}

	if len(messages) == 0 {
		return nil, ErrNoContent
	}

	return messages, nil
}
