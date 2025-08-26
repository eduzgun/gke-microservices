// internal/user/user_controller.go
package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/eduzgun/gke-microservices/internal/models"
)

type userControllerImpl struct {
	service UserService
	logger  *slog.Logger
}

type UserController interface {
	HandleGetUser(w http.ResponseWriter, r *http.Request)
	HandleCreateUser(w http.ResponseWriter, r *http.Request)
}

func NewUserController(service UserService, logger *slog.Logger) UserController {
	return &userControllerImpl{
		service: service,
		logger:  logger,
	}
}

func (c *userControllerImpl) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := c.service.GetUser(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		c.logger.Error("failed to get user", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (c *userControllerImpl) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := c.service.CreateUser(r.Context(), user); err != nil {
		c.logger.Error("failed to create user", "error", err)
		http.Error(w, "failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "created", "email": user.Email})
}
