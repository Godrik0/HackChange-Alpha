package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
)

func (h *ClientHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}

func (h *ClientHandler) respondError(w http.ResponseWriter, status int, message string) {
	h.respondJSON(w, status, dto.ErrorResponse{Error: message})
}
