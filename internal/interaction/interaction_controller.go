// internal/interaction/interaction_controller.go
package interaction

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/eduzgun/gke-microservices/internal/errs"
)

type interactionControllerImpl struct {
	service InteractionService
	logger  *slog.Logger
}

type InteractionController interface {
	HandleCreateComment(w http.ResponseWriter, r *http.Request)
	HandleToggleLike(w http.ResponseWriter, r *http.Request)
	HandleGetInteractions(w http.ResponseWriter, r *http.Request)
}

func NewInteractionController(service InteractionService, logger *slog.Logger) InteractionController {
	return &interactionControllerImpl{
		service: service,
		logger:  logger,
	}
}

func (c *interactionControllerImpl) HandleCreateComment(w http.ResponseWriter, r *http.Request) {
	// Add debugging
	c.logger.Info("HandleCreateComment called",
		"method", r.Method,
		"url", r.URL.Path,
		"full_url", r.URL.String())

	if r.Method != http.MethodPost {
		c.logger.Error("Wrong method", "method", r.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	philID, err := getIDFromPath(r, "/philosophers/")
	if err != nil {
		http.Error(w, "invalid philosopher ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// ✅ Get the created comment back from service
	comment, err := c.service.CreateComment(r.Context(), userID, philID, req.Content)
	if err != nil {
		if errors.Is(err, errs.ErrCommentContentRequired) {
			http.Error(w, "please add content to your comment before submitting", http.StatusBadRequest)
			return
		}
		c.logger.Error("failed to create comment", "error", err)
		http.Error(w, "failed to create comment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment) // ✅ Return the actual comment
}

func (c *interactionControllerImpl) HandleToggleLike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := getUserIDFromContext(r)
	if userID == 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	philID, err := getIDFromPath(r, "/philosophers/")
	if err != nil {
		http.Error(w, "invalid philosopher ID", http.StatusBadRequest)
		return
	}

	// Try to remove like first (assume it exists)
	err = c.service.RemoveLike(r.Context(), userID, philID)
	if err == nil {
		json.NewEncoder(w).Encode(map[string]bool{"liked": false})
		return
	}

	// If not found, create it
	if strings.Contains(err.Error(), "not found") {
		if err := c.service.CreateLike(r.Context(), userID, philID); err != nil {
			c.logger.Error("failed to like", "error", err)
			http.Error(w, "failed to like", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"liked": true})
		return
	}

	http.Error(w, "server error", http.StatusInternalServerError)
}

func (c *interactionControllerImpl) HandleGetInteractions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	philID, err := getIDFromPath(r, "/philosophers/")
	if err != nil {
		http.Error(w, "invalid philosopher ID", http.StatusBadRequest)
		return
	}

	interactions, err := c.service.GetInteractionsForPhilosopher(r.Context(), philID)
	if err != nil {
		c.logger.Error("failed to get interactions", "error", err)
		http.Error(w, "failed to load interactions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(interactions)
}

// Helper: extract ID from /philosophers/123/comments
func getIDFromPath(r *http.Request, prefix string) (int, error) {
	path := r.URL.Path
	start := strings.Index(path, prefix)
	if start == -1 {
		return 0, fmt.Errorf("invalid path")
	}
	start += len(prefix)
	end := strings.Index(path[start:], "/")
	if end == -1 {
		end = len(path)
	} else {
		end += start
	}
	return strconv.Atoi(path[start:end])
}

// Helper: get user ID from context (set by auth middleware)
func getUserIDFromContext(r *http.Request) int {
	if uid, ok := r.Context().Value("user_id").(int); ok {
		return uid
	}
	return 0
}
