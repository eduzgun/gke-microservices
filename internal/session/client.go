package session

import (
	"context"
	"fmt"
	"time"

	pb "github.com/eduzgun/gke-microservices/internal/proto/gen/go/session"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client wraps the gRPC session client
type Client struct {
	client pb.SessionServiceClient
	conn   *grpc.ClientConn
}

// NewClient creates a new session service client
func NewClient(sessionServiceAddr string) (*Client, error) {
	conn, err := grpc.NewClient(sessionServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to session service: %w", err)
	}

	client := pb.NewSessionServiceClient(conn)
	return &Client{client: client, conn: conn}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// CreateSession creates a new session
func (c *Client) CreateSession(ctx context.Context, userID int32) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.CreateSession(ctx, &pb.CreateSessionRequest{
		UserId: userID,
	})
	if err != nil {
		return "", fmt.Errorf("grpc call failed: %w", err)
	}

	if resp.Error != "" {
		return "", fmt.Errorf("session creation failed: %s", resp.Error)
	}

	return resp.SessionId, nil
}

// ValidateSession validates a session
func (c *Client) ValidateSession(ctx context.Context, sessionID string) (bool, int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.ValidateSession(ctx, &pb.ValidateSessionRequest{
		SessionId: sessionID,
	})
	if err != nil {
		return false, 0, fmt.Errorf("grpc call failed: %w", err)
	}

	if !resp.Valid {
		return false, 0, fmt.Errorf("session invalid: %s", resp.Error)
	}

	return resp.Valid, int(resp.UserId), nil
}

// InvalidateSession invalidates a session
func (c *Client) InvalidateSession(ctx context.Context, sessionID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := c.client.InvalidateSession(ctx, &pb.InvalidateSessionRequest{
		SessionId: sessionID,
	})
	if err != nil {
		return fmt.Errorf("grpc call failed: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("session invalidation failed: %s", resp.Error)
	}

	return nil
}
