package philosopher

import (
	"encoding/json"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/eduzgun/gke-microservices/internal/errs"
	"github.com/eduzgun/gke-microservices/internal/models"
)

type philosopherControllerImpl struct {
	service PhilosopherService
	logger  *slog.Logger
}

type PhilosopherController interface {
	HandleGetPhilosophers(w http.ResponseWriter, r *http.Request)
	HandleCreatePhilosopher(w http.ResponseWriter, r *http.Request)
	HandleGetPhilosopher(w http.ResponseWriter, r *http.Request)
}

func NewPhilosopherController(service PhilosopherService, logger *slog.Logger) PhilosopherController {
	return &philosopherControllerImpl{
		service: service,
		logger:  logger,
	}
}

func (pc *philosopherControllerImpl) HandleGetPhilosophers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	philosophers, err := pc.service.GetPhilosophers(r.Context())
	if err != nil {
		pc.logger.Error("failed to fetch philosophers", "error", err)
		http.Error(w, "failed to fetch philosophers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(philosophers); err != nil {
		// Can't use http.Error after WriteHeader
		return
	}
}

func (pc *philosopherControllerImpl) HandleGetPhilosopher(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract 'id' from URL path: /philosophers/3
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "invalid URL path", http.StatusBadRequest)
		return
	}

	idStr := parts[len(parts)-1] // Last part of /philosophers/3
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid philosopher ID", http.StatusBadRequest)
		return
	}

	philosopher, err := pc.service.GetPhilosopher(r.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			http.Error(w, "philosopher not found", http.StatusNotFound)
			return
		}
		pc.logger.Error("Failed to get philosopher", "error", err)
		http.Error(w, "failed to fetch philosopher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(philosopher); err != nil {
		pc.logger.Error("Failed to encode philosopher", "error", err)
		return
	}
}

func (pc *philosopherControllerImpl) HandleCreatePhilosopher(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var phil models.Philosopher

	if err := json.NewDecoder(r.Body).Decode(&phil); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := pc.service.CreatePhilosopher(r.Context(), phil); err != nil {
		pc.logger.Error("creating philosopher", "error", err)
		http.Error(w, "failed to create philosopher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "created",
		"name":   phil.Name,
	}); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}
