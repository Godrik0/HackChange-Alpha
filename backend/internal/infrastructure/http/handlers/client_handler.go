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
	clientService  interfaces.ClientService
	scoringService interfaces.ScoringService
	importService  interfaces.ImportService
	logger         interfaces.Logger
}

func NewClientHandler(clientService interfaces.ClientService, scoringService interfaces.ScoringService, importService interfaces.ImportService, logger interfaces.Logger) *ClientHandler {
	return &ClientHandler{
		clientService:  clientService,
		scoringService: scoringService,
		importService:  importService,
		logger:         logger.With("component", "ClientHandler"),
	}
}

// @Summary      Получение клиента
// @Description  Возвращает полные данные клиента по его ID
// @Tags         clients
// @Produce      json
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  dto.ClientResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/clients/{id} [get]
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

// @Summary      Поиск клиентов
// @Description  Поиск по ФИО и дате рождения (формат даты: DD-MM-YYYY). Требуется хотя бы один параметр.
// @Tags         clients
// @Produce      json
// @Param        first_name  query     string  false  "Имя (частичное совпадение)"
// @Param        last_name   query     string  false  "Фамилия (частичное совпадение)"
// @Param        middle_name query     string  false  "Отчество (частичное совпадение)"
// @Param        birth_date  query     string  false  "Дата рождения (DD-MM-YYYY)"
// @Success      200         {array}   dto.ClientResponse
// @Failure      400         {object}  dto.ErrorResponse
// @Failure      500         {object}  dto.ErrorResponse
// @Router       /api/clients/search [get]
func (h *ClientHandler) SearchClients(w http.ResponseWriter, r *http.Request) {
	params := dto.SearchParams{
		FirstName:  r.URL.Query().Get("first_name"),
		LastName:   r.URL.Query().Get("last_name"),
		MiddleName: r.URL.Query().Get("middle_name"),
		BirthDate:  r.URL.Query().Get("birth_date"),
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

// @Summary      Расчет ML-скоринга
// @Description  Запускает ML-модель для расчета скора клиента и получения рекомендаций
// @Tags         scoring
// @Produce      json
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  dto.ScoringResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/clients/{id}/scoring [get]
func (h *ClientHandler) CalculateScoring(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.logger.Warn("Invalid client ID", "id", idStr)
		h.respondError(w, http.StatusBadRequest, "invalid client ID")
		return
	}

	result, err := h.scoringService.CalculateScoring(r.Context(), id)
	if err != nil {
		if errors.Is(err, domainerrors.ErrClientNotFound) {
			h.respondError(w, http.StatusNotFound, "client not found")
			return
		}
		h.logger.Error("Failed to calculate scoring", "id", id, "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to calculate scoring")
		return
	}

	h.respondJSON(w, http.StatusOK, result)
}

// @Summary      Создание клиента
// @Description  Создает нового клиента с переданными данными (ФИО, дата рождения, признаки для ML)
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        input body dto.CreateClientRequest true "Данные клиента"
// @Success      201  {object}  dto.ClientResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/clients [post]
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

// UpdateClient обновляет данные клиента
// @Summary      Обновление клиента
// @Description  Обновляет данные существующего клиента. Все поля опциональны.
// @Tags         clients
// @Accept       json
// @Produce      json
// @Param        id    path   int                      true  "Client ID"
// @Param        input body   dto.UpdateClientRequest  true  "Данные для обновления"
// @Success      200  {object}  dto.ClientResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/clients/{id} [put]
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

// DeleteClient удаляет клиента
// @Summary      Удаление клиента
// @Description  Удаляет клиента из базы данных по его ID
// @Tags         clients
// @Produce      json
// @Param        id   path      int  true  "Client ID"
// @Success      200  {object}  dto.SuccessResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/clients/{id} [delete]
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

// ListClients возвращает список клиентов с пагинацией
// @Summary      Список клиентов
// @Description  Возвращает список всех клиентов с пагинацией (offset-based)
// @Tags         clients
// @Produce      json
// @Param        limit   query     int  false  "Количество записей (по умолчанию 100, макс 1000)"
// @Param        offset  query     int  false  "Смещение (по умолчанию 0)"
// @Success      200  {array}   dto.ClientResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /api/clients [get]
func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	limit := 100
	offset := 0

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		if parsedOffset, err := strconv.Atoi(offsetStr); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	clients, err := h.clientService.ListClients(r.Context(), limit, offset)
	if err != nil {
		h.logger.Error("Failed to list clients", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to list clients")
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

// ImportClientsCSV загружает клиентов из CSV файла
// @Summary      Импорт клиентов из CSV
// @Description  Загружает клиентов из CSV файла. Поддерживает две схемы: простую (first_name,last_name,birth_date,...) и полную (с features в JSON). Формат даты: DD-MM-YYYY или YYYY-MM-DD
// @Tags         clients
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "CSV файл с данными клиентов"
// @Success      200   {object}  interfaces.ImportStats  "Результат импорта с количеством успешных/неудачных записей"
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /api/clients/import [post]
func (h *ClientHandler) ImportClientsCSV(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Error("Failed to parse multipart form", "error", err)
		h.respondError(w, http.StatusBadRequest, "file too large or invalid form data")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to get file from form", "error", err)
		h.respondError(w, http.StatusBadRequest, "file is required")
		return
	}
	defer file.Close()

	stats, err := h.importService.ImportClientsCSV(r.Context(), file)
	if err != nil {
		h.logger.Error("Failed to import CSV", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to import clients")
		return
	}

	if len(stats.Errors) > 20 {
		stats.Errors = append(stats.Errors[:20], "... и другие ошибки")
	}

	h.logger.Info("CSV import completed", "success", stats.SuccessCount, "failures", stats.FailureCount)
	h.respondJSON(w, http.StatusOK, stats)
}
