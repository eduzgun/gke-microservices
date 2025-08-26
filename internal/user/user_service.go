// internal/user/user_service.go
package user

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/eduzgun/gke-microservices/internal/models"
)

type userServiceImpl struct {
	repo UserRepo
}

type UserService interface {
	GetUser(ctx context.Context, id int) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) error
	Authenticate(ctx context.Context, email, password string) (*models.User, error)
}

func NewUserService(repo UserRepo) UserService {
	return &userServiceImpl{repo: repo}
}

func (s *userServiceImpl) GetUser(ctx context.Context, id int) (*models.User, error) {
	user, err := s.repo.GetUser(ctx, int32(id))
	if err != nil {
		return nil, fmt.Errorf("getting user by id: %w", err)
	}
	return &user, nil
}

func (s *userServiceImpl) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("getting user by email: %w", err)
	}
	return &user, nil
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user models.User) error {
	if user.Username == "" || user.Email == "" || user.Password == "" {
		return fmt.Errorf("username, email, and password are required")
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	user.Password = string(hash)

	_, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	return nil
}

func (s *userServiceImpl) Authenticate(ctx context.Context, email, password string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Don't return password
	user.Password = ""
	return &user, nil
}
