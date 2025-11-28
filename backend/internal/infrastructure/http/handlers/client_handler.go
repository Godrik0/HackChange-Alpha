package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	domainerrors "github.com/Godrik0/HackChange-Alpha/backend/internal/domain/errors"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/go-chi/chi/v5"
)

type ClientHandler struct {
	clientService interfaces.ClientService
	logger        interfaces.Logger
}

func NewClientHandler(clientService interfaces.ClientService, logger interfaces.Logger) *ClientHandler {
	return &ClientHandler{
		clientService: clientService,
		logger:        logger.With("component", "ClientHandler"),
	}
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid client ID", "id", idStr)
		h.respondError(w, http.StatusBadRequest, "invalid client ID")
		return
	}

	client, err := h.clientService.GetClient(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			h.respondError(w, http.StatusNotFound, "client not found")
			return
		}
		h.logger.Error("Failed to get client", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to get client")
		return
	}

	response, err := dto.FromModel(client)
	if err != nil {
		h.logger.Error("Failed to convert model to DTO", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

func (h *ClientHandler) SearchClients(w http.ResponseWriter, r *http.Request) {
	params := dto.SearchParams{
		FirstName: r.URL.Query().Get("first_name"),
		LastName:  r.URL.Query().Get("last_name"),
		BirthDate: r.URL.Query().Get("birth_date"),
	}

	if params.IsEmpty() {
		h.respondError(w, http.StatusBadRequest, "at least one search parameter is required")
		return
	}

	clients, err := h.clientService.SearchClients(r.Context(), params)
	if err != nil {
		h.logger.Error("Failed to search clients", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to search clients")
		return
	}

	responses, err := dto.FromModels(clients)
	if err != nil {
		h.logger.Error("Failed to convert models to DTOs", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.respondJSON(w, http.StatusOK, responses)
}

func (h *ClientHandler) CalculateScoring(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid client ID", "id", idStr)
		h.respondError(w, http.StatusBadRequest, "invalid client ID")
		return
	}

	result, err := h.clientService.CalculateScoring(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			h.respondError(w, http.StatusNotFound, "client not found")
			return
		}
		h.logger.Error("Failed to calculate scoring", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to calculate scoring")
		return
	}

	response := &dto.ScoringResponse{
		Score:           result.Score,
		Recommendations: result.Recommendations,
		Factors:         result.Factors,
	}

	h.respondJSON(w, http.StatusOK, response)
}

func (h *ClientHandler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateClientRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	client, err := h.clientService.CreateClient(r.Context(), &req)
	if err != nil {
		h.logger.Error("Failed to create client", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to create client")
		return
	}

	response, err := dto.FromModel(client)
	if err != nil {
		h.logger.Error("Failed to convert model to DTO", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.respondJSON(w, http.StatusCreated, response)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid client ID", "id", idStr)
		h.respondError(w, http.StatusBadRequest, "invalid client ID")
		return
	}

	var req dto.UpdateClientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Warn("Invalid request body", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	client, err := h.clientService.UpdateClient(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			h.respondError(w, http.StatusNotFound, "client not found")
			return
		}
		h.logger.Error("Failed to update client", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to update client")
		return
	}

	response, err := dto.FromModel(client)
	if err != nil {
		h.logger.Error("Failed to convert model to DTO", "error", err)
		h.respondError(w, http.StatusInternalServerError, "internal error")
		return
	}

	h.respondJSON(w, http.StatusOK, response)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid client ID", "id", idStr)
		h.respondError(w, http.StatusBadRequest, "invalid client ID")
		return
	}

	if err := h.clientService.DeleteClient(r.Context(), id); err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			h.respondError(w, http.StatusNotFound, "client not found")
			return
		}
		h.logger.Error("Failed to delete client", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to delete client")
		return
	}

	h.respondJSON(w, http.StatusOK, dto.SuccessResponse{Message: "client deleted successfully"})
}
