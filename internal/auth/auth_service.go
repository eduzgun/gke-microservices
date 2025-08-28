package auth

import (
	"context"
	"fmt"

	"github.com/eduzgun/gke-microservices/internal/errs"
	"github.com/eduzgun/gke-microservices/internal/session"
	"github.com/eduzgun/gke-microservices/internal/user"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles HTTP authentication logic
type authServiceImpl struct {
	userService   user.UserService
	sessionClient *session.Client
}

type AuthService interface {
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context, sessionID string) error
}

func NewAuthService(userService user.UserService, sessionClient *session.Client) AuthService {
	return &authServiceImpl{
		userService:   userService,
		sessionClient: sessionClient,
	}
}

// Login handles user authentication and creates session
func (s *authServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	// Validate credentials
	user, err := s.userService.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("getting user by email: %w", err)
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errs.ErrInvalidCredentials
	}

	// Create session via gRPC
	sessionID, err := s.sessionClient.CreateSession(ctx, int32(user.ID))
	if err != nil {
		return "", fmt.Errorf("failed to create session with gRPC: %w", err)
	}

	return sessionID, nil
}

// Logout invalidates the session
func (s *authServiceImpl) Logout(ctx context.Context, sessionID string) error {
	return s.sessionClient.InvalidateSession(ctx, sessionID)
}
