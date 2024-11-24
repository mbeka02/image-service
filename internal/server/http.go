package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// APIError represents a structured error response
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// APIResponse represents a successful response
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ValidationError represents validation error details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

var ErrInvalidJSON = errors.New("invalid JSON payload")

// respondWithJSON handles writing JSON responses
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		return fmt.Errorf("failed to encode response: %w", err)
	}
	return nil
}

// respondWithError handles error responses in a consistent format
func respondWithError(w http.ResponseWriter, status int, err error) {
	apiError := APIError{
		Status:  status,
		Message: http.StatusText(status),
		Detail:  err.Error(),
	}

	respondWithJSON(w, status, apiError)
}

// parseJSON safely decodes JSON request bodies
func parseJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return fmt.Errorf("%w: %v", ErrInvalidJSON, err)
	}

	return nil
}

// validateRequest handles struct validation
func validateRequest(v interface{}) []ValidationError {
	if err := validate.Struct(v); err != nil {
		var validationErrors []ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Message: fmt.Sprintf("failed validation on '%s'", err.Tag()),
			})
		}
		return validationErrors
	}
	return nil
}
