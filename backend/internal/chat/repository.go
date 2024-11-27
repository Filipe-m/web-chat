package chat

import (
	"database/sql"
	"fmt"
)

type Repository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{DB: db}
}

func (repo *Repository) Create(chat *Chat, user int) error {

	query := "INSERT INTO chats (name, created_by) VALUES ($1, $2) RETURNING id, name, created_by, created_at, updated_at"

	err := repo.DB.QueryRow(query, chat.Name, user).Scan(&chat.ID, &chat.Name, &chat.Created_by, &chat.CreatedAt, &chat.UpdatedAt)

	if err != nil {
		return fmt.Errorf("erro ao criar chat: %v", err)
	}

	return nil
}

func (repo *Repository) GetChats() ([]Chat, error) {

	query := "SELECT id, name, created_by, created_at, updated_at FROM chats"

	var chats []Chat

	rows, err := repo.DB.Query(query)

	if err != nil {
		return nil, fmt.Errorf("erro ao executar a query: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.ID, &chat.Name, &chat.Created_by, &chat.CreatedAt, &chat.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("erro ao escanear os resultados: %w", err)
		}
		chats = append(chats, chat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro ao iterar os resultados: %w", err)
	}

	return chats, nil

}

func (repo *Repository) Delete(chatID, userID int) error {
	query := "DELETE FROM chats WHERE created_by = $1 AND id = $2"

	result, err := repo.DB.Exec(query, userID, chatID)
	if err != nil {
		return fmt.Errorf("erro ao executar a query de deleção: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("erro ao verificar linhas afetadas: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("nenhum registro encontrado para deletar")
	}

	return nil
}

func (repo *Repository) GetChatByID(chatID int, chat *Chat) error {

	query := "SELECT id, name, created_by, created_at, updated_at FROM chats where id = $1"

	err := repo.DB.QueryRow(query, chatID).Scan(&chat.ID, &chat.Name, &chat.Created_by, &chat.CreatedAt, &chat.UpdatedAt)

	if err != nil {
		return fmt.Errorf("chat não encontrado: %v", err)
	}

	return nil
}
