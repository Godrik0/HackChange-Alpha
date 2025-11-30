package handlers

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/dto"
	domainerrors "github.com/Godrik0/HackChange-Alpha/backend/internal/domain/errors"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/go-chi/chi/v5"
)

type ClientHandler struct {
	clientService  interfaces.ClientService
	scoringService interfaces.ScoringService
	logger         interfaces.Logger
}

func NewClientHandler(clientService interfaces.ClientService, scoringService interfaces.ScoringService, logger interfaces.Logger) *ClientHandler {
	return &ClientHandler{
		clientService:  clientService,
		scoringService: scoringService,
		logger:         logger.With("component", "ClientHandler"),
	}
}

// GetClient получает клиента по ID
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

// SearchClients ищет клиентов
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

// CalculateScoring рассчитывает скоринг клиента
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

// CreateClient создает нового клиента
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
// @Success      200   {object}  map[string]interface{}  "Результат импорта с количеством успешных/неудачных записей"
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /api/clients/import [post]
func (h *ClientHandler) ImportClientsCSV(w http.ResponseWriter, r *http.Request) {
	// Ограничение размера файла (10MB)
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

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	// Читаем заголовки
	headers, err := reader.Read()
	if err != nil {
		h.logger.Error("Failed to read CSV headers", "error", err)
		h.respondError(w, http.StatusBadRequest, "invalid CSV format")
		return
	}

	// Нормализуем заголовки (trim, lowercase)
	for i := range headers {
		headers[i] = strings.TrimSpace(strings.ToLower(headers[i]))
	}

	successCount := 0
	failureCount := 0
	var errors []string

	lineNum := 1 // Начинаем с 1 (заголовки)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		lineNum++

		if err != nil {
			h.logger.Warn("Failed to read CSV line", "line", lineNum, "error", err)
			errors = append(errors, fmt.Sprintf("Line %d: %v", lineNum, err))
			failureCount++
			continue
		}

		if len(record) != len(headers) {
			errors = append(errors, fmt.Sprintf("Line %d: column count mismatch", lineNum))
			failureCount++
			continue
		}

		// Создаем map из записи
		rowData := make(map[string]string)
		for i, header := range headers {
			rowData[header] = strings.TrimSpace(record[i])
		}

		// Парсим клиента из строки
		clientReq, err := h.parseClientFromCSVRow(rowData, headers)
		if err != nil {
			errors = append(errors, fmt.Sprintf("Line %d: %v", lineNum, err))
			failureCount++
			continue
		}

		// Создаем клиента
		if _, err := h.clientService.CreateClient(r.Context(), clientReq); err != nil {
			h.logger.Warn("Failed to create client from CSV", "line", lineNum, "error", err)
			errors = append(errors, fmt.Sprintf("Line %d: %v", lineNum, err))
			failureCount++
			continue
		}

		successCount++
	}

	result := map[string]interface{}{
		"success_count": successCount,
		"failure_count": failureCount,
		"total":         successCount + failureCount,
	}

	if len(errors) > 0 {
		// Ограничиваем количество ошибок в ответе
		if len(errors) > 20 {
			result["errors"] = append(errors[:20], "... и другие ошибки")
		} else {
			result["errors"] = errors
		}
	}

	h.logger.Info("CSV import completed", "success", successCount, "failures", failureCount)
	h.respondJSON(w, http.StatusOK, result)
}

// parseClientFromCSVRow парсит клиента из CSV строки
func (h *ClientHandler) parseClientFromCSVRow(row map[string]string, headers []string) (*dto.CreateClientRequest, error) {
	req := &dto.CreateClientRequest{}

	// Обязательные поля
	firstName, ok := row["first_name"]
	if !ok || firstName == "" {
		return nil, errors.New("first_name is required")
	}
	req.FirstName = firstName

	lastName, ok := row["last_name"]
	if !ok || lastName == "" {
		return nil, errors.New("last_name is required")
	}
	req.LastName = lastName

	// Парсим дату рождения (поддерживаем форматы DD-MM-YYYY и YYYY-MM-DD)
	birthDateStr, ok := row["birth_date"]
	if !ok || birthDateStr == "" {
		return nil, errors.New("birth_date is required")
	}

	var birthDate time.Time
	var parseErr error

	// Пробуем разные форматы
	formats := []string{"02-01-2006", "2006-01-02", "02/01/2006", "2006/01/02"}
	for _, format := range formats {
		birthDate, parseErr = time.Parse(format, birthDateStr)
		if parseErr == nil {
			break
		}
	}

	if parseErr != nil {
		return nil, fmt.Errorf("invalid birth_date format: %s (expected DD-MM-YYYY or YYYY-MM-DD)", birthDateStr)
	}

	req.BirthDate = birthDate.Format(dto.DateFormat)

	// Опциональные поля
	if middleName, ok := row["middle_name"]; ok {
		req.MiddleName = middleName
	}

	// Собираем все дополнительные поля как features
	features := make(map[string]interface{})

	// Список базовых полей, которые не должны попасть в features
	baseFields := map[string]bool{
		"first_name":  true,
		"last_name":   true,
		"middle_name": true,
		"birth_date":  true,
		"phone":       true,
		"email":       true,
		"address":     true,
	}

	for _, header := range headers {
		if baseFields[header] {
			continue
		}

		value := row[header]
		if value == "" {
			continue
		}

		// Пытаемся распарсить как число
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			features[header] = floatVal
		} else {
			features[header] = value
		}
	}

	if len(features) > 0 {
		req.Features = features
	}

	return req, nil
}
