package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"github.com/mbeka02/image-service/internal/auth"
	"github.com/mbeka02/image-service/internal/database"
	"github.com/mbeka02/image-service/internal/models"
)

var validate *validator.Validate

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	params := models.CreateUserRequest{}

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

	passwordHash, err := auth.HashPassword(params.Password)
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
	params := models.LoginRequest{}
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
	err = auth.ComparePassword(params.Password, user.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err)
		return
	}
	userResponse := models.NewUserResponse(user)
	token, err := s.AuthMaker.Create(user.Email, user.UserID, s.AccessTokenDuration)
	resp := models.LoginResponse{
		AccessToken: token,
		User:        userResponse,
	}
	if err := respondWithJSON(w, http.StatusCreated, resp); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

// func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
// 	users, err := s.Store.GetUsers(r.Context(), database.GetUsersParams{
// 		Limit:  10,
// 		Offset: 0,
// 	})
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err)
// 		return
// 	}
// 	response := APIResponse{
// 		Status:  http.StatusOK,
// 		Data:    users,
// 		Message: "users:",
// 	}
// 	respondWithJSON(w, http.StatusOK, response)
// 	return
// }

func (s *Server) handleImageUpload(w http.ResponseWriter, r *http.Request) {
	// get the file
	_, fileHeader, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Errorf("bad request:%v", err))
		return
	}
	// upload the file to GC storage
	uploadResponse, err := s.FileStorage.Upload(r.Context(), fileHeader)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Errorf("internal server error : %v", err))
		return
	}
	payload, err := getAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	// save to DB
	createdImage, err := s.Store.CreateImage(r.Context(), database.CreateImageParams{
		UserID:     payload.UserID,
		FileName:   uploadResponse.FileName,
		StorageUrl: uploadResponse.StorageUrl,
		FileSize:   uploadResponse.Size,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	response := APIResponse{
		Status:  http.StatusOK,
		Data:    createdImage,
		Message: "uploaded",
	}
	respondWithJSON(w, http.StatusOK, response)
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
