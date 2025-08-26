// internal/user/user_repo.go
package user

import (
	"context"
	"fmt"

	"github.com/eduzgun/gke-microservices/internal/db"
	"github.com/eduzgun/gke-microservices/internal/models"
)

type userRepoImpl struct {
	queries db.Querier
}

type UserRepo interface {
	GetUser(ctx context.Context, id int) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (int, error)
}

func NewUserRepo(queries db.Querier) UserRepo {
	return &userRepoImpl{queries: queries}
}

func (r *userRepoImpl) GetUser(ctx context.Context, id int) (models.User, error) {
	dbUser, err := r.queries.GetUser(ctx, int32(id))
	if err != nil {
		return models.User{}, fmt.Errorf("querier: getting user: %w", err)
	}

	return models.User{
		ID:       int(dbUser.ID),
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Password: dbUser.Password,
	}, nil
}

func (r *userRepoImpl) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, fmt.Errorf("querier: getting user by email: %w", err)
	}

	return models.User{
		ID:       int(dbUser.ID),
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Password: dbUser.Password,
	}, nil
}

func (r *userRepoImpl) CreateUser(ctx context.Context, user models.User) (int, error) {
	params := db.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Password: user.Password,
	}

	id, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("querier: creating user: %w", err)
	}

	return int(id), nil
}
