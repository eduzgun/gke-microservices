// internal/interaction/interaction_service.go
package interaction

import (
	"context"
	"fmt"

	"github.com/eduzgun/gke-microservices/internal/models"
)

type interactionServiceImpl struct {
	repo InteractionRepo
}

type InteractionService interface {
	CreateComment(ctx context.Context, userID, philID int, content string) error
	CreateLike(ctx context.Context, userID, philID int) error
	RemoveLike(ctx context.Context, userID, philID int) error
	GetInteractionsForPhilosopher(ctx context.Context, philID int) ([]models.InteractionResponse, error)
}

func NewInteractionService(repo InteractionRepo) InteractionService {
	return &interactionServiceImpl{repo: repo}
}

func (s *interactionServiceImpl) CreateComment(ctx context.Context, userID, philID int, content string) error {
	if content == "" {
		return fmt.Errorf("comment content is required")
	}

	interaction := models.Interaction{
		UserID:        userID,
		PhilosopherID: philID,
		Type:          "comment",
		Content:       content,
	}

	return s.repo.CreateInteraction(ctx, interaction)
}

func (s *interactionServiceImpl) CreateLike(ctx context.Context, userID, philID int) error {
	// Prevent duplicate likes
	_, err := s.repo.GetUserInteraction(ctx, userID, philID, "like")
	if err == nil {
		return fmt.Errorf("already liked")
	}

	interaction := models.Interaction{
		UserID:        userID,
		PhilosopherID: philID,
		Type:          "like",
	}

	return s.repo.CreateInteraction(ctx, interaction)
}

func (s *interactionServiceImpl) RemoveLike(ctx context.Context, userID, philID int) error {
	interaction, err := s.repo.GetUserInteraction(ctx, userID, philID, "like")
	if err != nil {
		return fmt.Errorf("like not found")
	}

	return s.repo.DeleteInteraction(ctx, interaction.ID)
}

func (s *interactionServiceImpl) GetInteractionsForPhilosopher(ctx context.Context, philID int) ([]models.InteractionResponse, error) {
	dbInteractions, err := s.repo.GetInteractionsByPhilosopher(ctx, philID)
	if err != nil {
		return nil, err
	}

	// You'd join with users table to get usernames
	// For now, return basic mapping
	var responses []models.InteractionResponse
	for _, i := range dbInteractions {
		responses = append(responses, models.InteractionResponse{
			ID:        int(i.ID),
			Username:  fmt.Sprintf("user-%d", i.UserID), // Replace with real username later
			Type:      i.Type,
			Content:   i.Content,
			CreatedAt: i.CreatedAt,
		})
	}

	return responses, nil
}
