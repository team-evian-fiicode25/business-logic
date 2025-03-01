package handlers

import (
	"business_logic/internal/models"
	"encoding/json"
	"net/http"
)

type AuthoritiesAuthHandler struct{}

func (h *AuthoritiesAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request models.AuthRequest
	var response models.AuthResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response = models.AuthResponse{
		Session_id: "test_session",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *AuthoritiesAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request models.AuthRequest
	var response models.AuthResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response = models.AuthResponse{
		Session_id: "test_session",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
