package interaction

import (
	"context"
	"fmt"

	"github.com/eduzgun/gke-microservices/internal/errs"
	"github.com/eduzgun/gke-microservices/internal/models"
)

type interactionServiceImpl struct {
	repo InteractionRepo
}

type InteractionService interface {
	CreateComment(ctx context.Context, userID, philID int, username string, content string) (*models.Comment, error)
	CreateLike(ctx context.Context, userID, philID int, username string) error
	RemoveLike(ctx context.Context, userID, philID int, username string) error
	GetInteractionsForPhilosopher(ctx context.Context, philID int) ([]models.InteractionResponse, error)
}

func NewInteractionService(repo InteractionRepo) InteractionService {
	return &interactionServiceImpl{repo: repo}
}

func (s *interactionServiceImpl) CreateComment(
	ctx context.Context,
	userID, philID int,
	username string,
	content string,
) (*models.Comment, error) {
	if content == "" {
		return nil, errs.ErrCommentContentRequired
	}

	dbComment, err := s.repo.CreateCommentInteraction(ctx, userID, philID, username, content)
	if err != nil {
		return nil, fmt.Errorf("creating comment interaction: %w", err)
	}

	comment := &models.Comment{
		ID:            int(dbComment.ID),
		UserID:        int(dbComment.UserID),
		PhilosopherID: int(dbComment.PhilosopherID),
		Username:      dbComment.Username,
		Content:       dbComment.Content,
		CreatedAt:     dbComment.CreatedAt,
	}

	return comment, nil
}

func (s *interactionServiceImpl) CreateLike(ctx context.Context, userID, philID int, username string) error {
	_, err := s.repo.GetUserInteraction(ctx, userID, philID, "like")
	if err == nil {
		return fmt.Errorf("already liked")
	}

	err = s.repo.CreateLikeInteraction(ctx, userID, philID, username)
	if err != nil {
		return fmt.Errorf("creating like interaction: %w", err)
	}
	return nil
}

func (s *interactionServiceImpl) RemoveLike(ctx context.Context, userID, philID int, username string) error {
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

	var responses []models.InteractionResponse
	for _, i := range dbInteractions {
		responses = append(responses, models.InteractionResponse{
			ID:        int(i.ID),
			Username:  fmt.Sprintf("user-%d", i.UserID),
			Type:      i.Type,
			Content:   i.Content,
			CreatedAt: i.CreatedAt,
		})
	}

	return responses, nil
}
