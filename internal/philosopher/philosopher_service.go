package philosopher

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/eduzgun/gke-microservices/internal/errs"
	"github.com/eduzgun/gke-microservices/internal/models"
)

type philosopherServiceImpl struct {
	repo PhilosopherRepo
}

type PhilosopherService interface {
	GetPhilosophers(ctx context.Context) ([]models.Philosopher, error)
	CreatePhilosopher(ctx context.Context, philosopher models.Philosopher) error
	GetPhilosopher(ctx context.Context, id int) (*models.Philosopher, error)
}

func NewPhilosopherService(repo PhilosopherRepo) PhilosopherService {
	return &philosopherServiceImpl{
		repo: repo,
	}
}

func (ps *philosopherServiceImpl) GetPhilosophers(ctx context.Context) ([]models.Philosopher, error) {
	return ps.repo.GetPhilosophers(ctx)
}

func (ps *philosopherServiceImpl) GetPhilosopher(ctx context.Context, id int) (*models.Philosopher, error) {

	philosopher, err := ps.repo.GetPhilosopher(ctx, int32(id))
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.ErrNotFound
		}

		return nil, fmt.Errorf("failed to fetch philosopher: %w", err)
	}

	return &philosopher, nil
}

func (ps *philosopherServiceImpl) CreatePhilosopher(ctx context.Context, philosopher models.Philosopher) error {
	if philosopher.Name == "" {
		return fmt.Errorf("name is required")
	}

	_, err := ps.repo.CreatePhilosopher(ctx, philosopher)
	if err != nil {
		return fmt.Errorf("creating philosopher in repo: %w", err)
	}

	return nil
}
