package philosopher

import (
	"encoding/json"
	"net/http"
)

type philosopherControllerImpl struct {
	service PhilosopherService
}

type PhilosopherController interface {
	HandleGetPhilosophers(w http.ResponseWriter, r *http.Request)
}

func NewPhilosopherController(service PhilosopherService) PhilosopherController {
	return &philosopherControllerImpl{
		service: service,
	}
}

func (pc *philosopherControllerImpl) HandleGetPhilosophers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	philosophers, err := pc.service.GetPhilosophers(r.Context())
	if err != nil {
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
