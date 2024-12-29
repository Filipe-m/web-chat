package rooms

import (
	"time"
)

type Room struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	CreatedBy int       `json:"created_by"`
	RoomId    int       `json:"room_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
