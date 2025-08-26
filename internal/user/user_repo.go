// internal/user/user_repo.go
package user

import (
	"context"
	"fmt"

	"github.com/eduzgun/gke-microservices/internal/db"
	"github.com/eduzgun/gke-microservices/internal/models"
	"github.com/jackc/pgx/v5/pgtype"
)

type userRepoImpl struct {
	queries db.Querier
}

type UserRepo interface {
	GetUser(ctx context.Context, id int32) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (int, error)
}

func NewUserRepo(queries db.Querier) UserRepo {
	return &userRepoImpl{queries: queries}
}

func (r *userRepoImpl) GetUser(ctx context.Context, id int32) (models.User, error) {
	dbUser, err := r.queries.GetUser(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("querier: getting user: %w", err)
	}

	return models.User{
		ID:       int(dbUser.ID),
		Username: dbUser.Username,
		Email:    dbUser.Email,
		Password: pgTextToString(dbUser.Password),
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
		Password: pgTextToString(dbUser.Password),
	}, nil
}

func (r *userRepoImpl) CreateUser(ctx context.Context, user models.User) (int, error) {
	params := db.CreateUserParams{
		Username: user.Username,
		Email:    user.Email,
		Password: stringToPgText(user.Password),
	}

	id, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("querier: creating user: %w", err)
	}

	return int(id), nil
}

// Helper: pgtype.Text → string
func pgTextToString(t pgtype.Text) string {
	if t.Valid {
		return t.String
	}
	return ""
}

// Helper: string → pgtype.Text
func stringToPgText(s string) pgtype.Text {
	if s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: s, Valid: true}
}
