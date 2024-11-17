package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
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

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				respondWithError(w, http.StatusForbidden, errors.New("forbidden: the username or email are already in use"))
				return
			}
		}
		respondWithError(w, http.StatusInternalServerError, errors.New("failed to create user"))
		return
	}
	response := APIResponse{
		Status:  http.StatusCreated,
		Message: "User created successfully",
		Data:    user,
	}
	go func() {
		if err := s.Mailer.SendEmail(); err != nil {
			fmt.Println("unable to send the email:", err)

			return
		}
	}()

	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	params := LoginRequest{}
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
	user, err := s.Store.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, errors.New("unable to find user"))
		return
	}
	err = ComparePassword(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}
	userResponse := newUserResponse(user)
	token, err := s.AuthMaker.Create(user.Email, s.AccessTokenDuration)
	resp := LoginResponse{
		AccessToken: token,
		User:        userResponse,
	}
	if err := respondWithJSON(w, http.StatusCreated, resp); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.Store.GetUsers(r.Context(), database.GetUsersParams{
		Limit:  10,
		Offset: 0,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	response := APIResponse{
		Status:  http.StatusOK,
		Data:    users,
		Message: "users:",
	}
	respondWithJSON(w, http.StatusOK, response)
	return
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
