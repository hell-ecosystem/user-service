package model

import "time"

// DTO для клиента
type User struct {
	ID         string    `json:"id"`
	Email      *string   `json:"email,omitempty"`
	TelegramID *int64    `json:"telegram_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}
