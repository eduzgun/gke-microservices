package models

import "time"

type Interaction struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	PhilosopherID int    `json:"philosopher_id"`
	Username      string `json:"username"`
	Type          string `json:"type"` // "like" or "comment"
	Content       string `json:"content,omitempty"`
	CreatedAt     string `json:"created_at"`
}

// For responses
type InteractionResponse struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Type      string `json:"type"`
	Content   string `json:"content,omitempty"`
	CreatedAt string `json:"created_at"`
}

type Comment struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"userId" db:"user_id"`
	PhilosopherID int       `json:"philosopherId" db:"philosopher_id"`
	Username      string    `json:"username" db:"username"`
	Content       string    `json:"content" db:"content"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}
