// internal/interaction/interaction_repo.go
package interaction

import (
	"context"
	"fmt"

	"github.com/eduzgun/gke-microservices/internal/db"
	"github.com/eduzgun/gke-microservices/internal/models"
)

type interactionRepoImpl struct {
	queries db.Querier
}

type InteractionRepo interface {
	GetInteractionsByPhilosopher(ctx context.Context, philID int) ([]models.Interaction, error)
	CreateInteraction(ctx context.Context, interaction models.Interaction) error
	DeleteInteraction(ctx context.Context, id int) error
	GetUserInteraction(ctx context.Context, userID, philID int, typ string) (models.Interaction, error)
}

func NewInteractionRepo(queries db.Querier) InteractionRepo {
	return &interactionRepoImpl{queries: queries}
}

func (r *interactionRepoImpl) GetInteractionsByPhilosopher(ctx context.Context, philID int) ([]models.Interaction, error) {
	rows, err := r.queries.GetInteractionsByPhilosopher(ctx, int32(philID))
	if err != nil {
		return nil, fmt.Errorf("failed to get interactions: %w", err)
	}

	// Map: []db.GetInteractionsByPhilosopherRow → []models.Interaction
	interactions := make([]models.Interaction, 0, len(rows))
	for _, row := range rows {
		interactions = append(interactions, models.Interaction{
			ID:            int(row.ID),
			UserID:        int(row.UserID),
			PhilosopherID: int(row.PhilosopherID),
			Type:          row.Type,
			Content:       row.Content,
			CreatedAt:     row.CreatedAt.Format("2006-01-02"),
		})
	}

	return interactions, nil
}

func (r *interactionRepoImpl) CreateInteraction(ctx context.Context, interaction models.Interaction) error {
	params := db.CreateInteractionParams{
		UserID:        int32(interaction.UserID),
		PhilosopherID: int32(interaction.PhilosopherID),
		Type:          interaction.Type,
		Content:       interaction.Content,
	}

	return r.queries.CreateInteraction(ctx, params)
}

func (r *interactionRepoImpl) DeleteInteraction(ctx context.Context, id int) error {
	return r.queries.DeleteInteraction(ctx, int32(id))
}

func (r *interactionRepoImpl) GetUserInteraction(ctx context.Context, userID, philID int, typ string) (models.Interaction, error) {
	params := db.GetUserInteractionParams{
		UserID:        int32(userID),
		PhilosopherID: int32(philID),
		Type:          typ,
	}

	row, err := r.queries.GetUserInteraction(ctx, params)
	if err != nil {
		return models.Interaction{}, err
	}

	return models.Interaction{
		ID:            int(row.ID),
		UserID:        int(row.UserID),
		PhilosopherID: int(row.PhilosopherID),
		Type:          row.Type,
		Content:       row.Content,
		CreatedAt:     row.CreatedAt.Format("2006-01-02"),
	}, nil
}
