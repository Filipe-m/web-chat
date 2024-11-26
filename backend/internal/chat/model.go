package chat

import (
	"time"
)

type Chat struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Created_by int       `json:"created_by"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
