package models

type Interaction struct {
	ID            int    `json:"id"`
	UserID        int    `json:"user_id"`
	PhilosopherID int    `json:"philosopher_id"`
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
