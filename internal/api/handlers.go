package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
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
	if err := respondWithJSON(w, http.StatusOK, resp); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

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

func (s *Server) handleGetImages(w http.ResponseWriter, r *http.Request) {
	payload, err := getAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}

	// Get limit from query parameter, default to 10 if not provided
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10 // default limit
	}

	// Get offset from query parameter, default to 0 if not provided
	offsetStr := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0 // default offset
	}
	data, err := s.Store.GetUserImages(r.Context(), database.GetUserImagesParams{
		UserID: payload.UserID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	response := APIResponse{
		Status:  http.StatusOK,
		Data:    data,
		Message: "images",
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (s *Server) handleDeleteImage(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "imageId")
	imageId, err := strconv.Atoi(idParam)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, errors.New("invalid url param"))
		return
	}
	payload, err := getAuthPayload(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
	image, err := s.Store.GetImage(r.Context(), int64(imageId))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to get image"))
		return
	}
	if image.UserID != payload.UserID {
		respondWithError(w, http.StatusUnauthorized, errors.New("unauthorized!"))
		return
	}
	if err = s.FileStorage.Delete(r.Context(), image.FileName); err != nil {
		respondWithError(w, http.StatusInternalServerError, errors.New("unable to delete the image"))
		return
	}

	s.Store.DeleteUserImage(r.Context(), database.DeleteUserImageParams{
		UserID:  payload.UserID,
		ImageID: int64(imageId),
	})
	response := APIResponse{
		Status:  http.StatusOK,
		Message: "deleted the image sucessfully",
		Data:    nil,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}
