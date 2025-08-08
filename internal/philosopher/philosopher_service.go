package philosopher

import (
	"context"

	"github.com/eduzgun/gke-microservices/internal/models"
)

type philosopherServiceImpl struct {
	repo PhilosopherRepo
}

type PhilosopherService interface {
	GetPhilosophers(ctx context.Context) ([]models.Philosopher, error)
}

func NewPhilosopherService(repo PhilosopherRepo) PhilosopherService {
	return &philosopherServiceImpl{
		repo: repo,
	}
}

func (ps *philosopherServiceImpl) GetPhilosophers(ctx context.Context) ([]models.Philosopher, error) {
	return ps.repo.GetPhilosophers(ctx)
}
