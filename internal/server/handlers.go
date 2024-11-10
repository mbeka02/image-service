// package server
//
// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
//
// 	"github.com/go-playground/validator/v10"
// 	"github.com/mbeka02/image-service/internal/database"
// )
//
// var validate *validator.Validate
//
// // handle JSON responses
// func ToJson(w http.ResponseWriter, statusCode int, payload interface{}) error {
// 	w.Header().Add("Content-Type", "application/json")
// 	w.WriteHeader(statusCode)
//
// 	return json.NewEncoder(w).Encode(&payload)
// }
//
// func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
// 	params := CreateUserRequest{}
// 	json.NewDecoder(r.Body).Decode(&params)
// 	fmt.Println(params)
// 	validate := validator.New()
// 	if err := validate.Struct(params); err != nil {
// 		ToJson(w, http.StatusBadRequest, err.Error())
// 		return
// 	}
// 	passwordHash, err := HashPassword(params.Password)
// 	if err != nil {
// 		ToJson(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	user, err := s.Store.CreateUser(r.Context(), database.CreateUserParams{
// 		FullName: params.Fullname,
// 		Email:    params.Email,
// 		Password: passwordHash,
// 	})
// 	if err != nil {
// 		ToJson(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	ToJson(w, http.StatusCreated, user)
// }

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mbeka02/image-service/internal/database"
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

var (
	validate *validator.Validate

	ErrInvalidJSON    = errors.New("invalid JSON payload")
	ErrInvalidRequest = errors.New("invalid request")
)

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

	// TODO : consider adding structured error logs

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

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	params := CreateUserRequest{}

	if err := parseJSON(r, &params); err != nil {
		respondWithError(w, http.StatusBadRequest, err)
		return
	}

	if validationErrors := validateRequest(params); validationErrors != nil {
		respondWithJSON(w, http.StatusBadRequest, APIError{
			Status:  http.StatusBadRequest,
			Message: "Validation failed",
			Detail:  fmt.Sprintf("%v", validationErrors),
		})
		return
	}

	passwordHash, err := HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("failed to process password"))
		return
	}

	user, err := s.Store.CreateUser(r.Context(), database.CreateUserParams{
		FullName: params.Fullname,
		Email:    params.Email,
		Password: passwordHash,
	})
	if err != nil {
		// Here you might want to check for specific errors like unique constraint violations
		respondWithError(w, http.StatusInternalServerError, errors.New("failed to create user"))
		return
	}

	response := APIResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	}

	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
