package handlers

import (
	"business_logic/internal/models"
	"business_logic/internal/services"
	"context"
	"encoding/json"
	"net/http"
)

type ConsumerAuthHandler struct {
	consumerService *services.ConsumerService
}

func NewConsumerAuthHandler(client *services.ConsumerService) *ConsumerAuthHandler {
	return &ConsumerAuthHandler{
		consumerService: client,
	}
}

func (h *ConsumerAuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request models.AuthRequest
	var response models.AuthResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	newConsumer, err := h.consumerService.CreateConsumer(context.Background(), request.Username, request.Email, request.Phone_number, request.Password)
	if err != nil {
		response.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response = models.AuthResponse{
		Id:         newConsumer.NewLogin.GetId(),
		Session_id: newConsumer.NewLogin.GetId(),
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *ConsumerAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var request models.AuthRequest
	var response models.AuthResponse

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		response.Error = err.Error()

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	loginSession, err := h.consumerService.LogInWithPassword(context.Background(), request.Identifier, request.Password)
	if err != nil {
		response.Error = err.Error()

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	response = models.AuthResponse{
		Id:         loginSession.LoginSession.GetId(),
		Session_id: loginSession.LoginSession.GetIdentifyingToken(),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
