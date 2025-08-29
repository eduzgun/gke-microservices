package session

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	pb "github.com/eduzgun/gke-microservices/internal/proto/gen/go/session"
)

type SessionStore interface {
	Create(userID int32, sessionID string) error
	Validate(sessionID string) (int32, bool)
	Invalidate(sessionID string) bool
}

type InMemorySessionStore struct {
	sessions map[string]SessionEntry
}

type SessionEntry struct {
	UserID    int32
	ExpiresAt time.Time
}

func NewInMemorySessionStore() *InMemorySessionStore {
	return &InMemorySessionStore{
		sessions: make(map[string]SessionEntry),
	}
}

func (s *InMemorySessionStore) Create(userID int32, sessionID string) error {
	s.sessions[sessionID] = SessionEntry{
		UserID:    userID,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return nil
}

func (s *InMemorySessionStore) Validate(sessionID string) (int32, bool) {
	entry, exists := s.sessions[sessionID]
	if !exists || entry.ExpiresAt.Before(time.Now()) {
		return 0, false
	}
	return entry.UserID, true
}

func (s *InMemorySessionStore) Invalidate(sessionID string) bool {
	_, exists := s.sessions[sessionID]
	delete(s.sessions, sessionID)
	return exists
}

// GRPCSessionServer implements the gRPC SessionService
type GRPCSessionServer struct {
	pb.UnimplementedSessionServiceServer
	store SessionStore
	log   *slog.Logger
}

func NewGRPCSessionServer(store SessionStore, log *slog.Logger) *GRPCSessionServer {
	return &GRPCSessionServer{store: store, log: log}
}

func generateSessionID() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func (s *GRPCSessionServer) CreateSession(ctx context.Context, req *pb.CreateSessionRequest) (*pb.CreateSessionResponse, error) {
	sessionID := generateSessionID()
	if err := s.store.Create(req.UserId, sessionID); err != nil {
		return &pb.CreateSessionResponse{
			SessionId: "",
			Error:     "failed to create session",
		}, nil
	}

	s.log.Info("Session created", "user_id", req.UserId, "session_id", sessionID)
	return &pb.CreateSessionResponse{
		SessionId: sessionID,
		Error:     "",
	}, nil
}

func (s *GRPCSessionServer) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	userId, valid := s.store.Validate(req.SessionId)
	if !valid {
		return &pb.ValidateSessionResponse{
			Valid:  false,
			UserId: 0,
			Error:  "invalid or expired session",
		}, nil
	}

	return &pb.ValidateSessionResponse{
		Valid:  true,
		UserId: userId,
		Error:  "",
	}, nil
}

func (s *GRPCSessionServer) InvalidateSession(ctx context.Context, req *pb.InvalidateSessionRequest) (*pb.InvalidateSessionResponse, error) {
	success := s.store.Invalidate(req.SessionId)
	if !success {
		return &pb.InvalidateSessionResponse{
			Success: false,
			Error:   "session not found",
		}, nil
	}

	return &pb.InvalidateSessionResponse{
		Success: true,
		Error:   "",
	}, nil
}
