package errors

import "errors"

// Ошибки репозитория клиентов
var (
	ErrClientNotFound = errors.New("client not found")

	ErrInvalidClientID = errors.New("invalid client ID")

	ErrClientAlreadyExists = errors.New("client already exists")
)

// Ошибки валидации
var (
	ErrInvalidInput = errors.New("invalid input data")

	ErrRequiredField = errors.New("required field is missing")

	ErrInvalidFormat = errors.New("invalid data format")
)

// Ошибки ML сервиса
var (
	ErrMLServiceUnavailable = errors.New("ML service is unavailable")

	ErrMLPredictionFailed = errors.New("ML prediction failed")

	ErrInvalidFeatures = errors.New("invalid features for ML model")
)

// Ошибки базы данных
var (
	ErrDatabaseConnection = errors.New("database connection error")

	ErrDatabaseQuery = errors.New("database query error")

	ErrTransactionFailed = errors.New("transaction failed")
)
