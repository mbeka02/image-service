package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/lib/pq"
	"github.com/mbeka02/image-service/internal/auth"
	"github.com/mbeka02/image-service/internal/database"
	"github.com/mbeka02/image-service/internal/mailer"
	"github.com/mbeka02/image-service/internal/models"
)

type UserHandler struct {
	Store               *database.Store
	AuthMaker           auth.Maker
	Mailer              *mailer.Mailer
	AccessTokenDuration time.Duration
}

func (uh *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
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

	user, err := uh.Store.CreateUser(r.Context(), database.CreateUserParams{
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
		if err := uh.Mailer.SendEmail(); err != nil {
			fmt.Println("unable to send the email:", err)

			return
		}
	}()

	if err := respondWithJSON(w, http.StatusCreated, response); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}

func (uh *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
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
	user, err := uh.Store.GetUserByEmail(r.Context(), params.Email)
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
	token, err := uh.AuthMaker.Create(user.Email, user.UserID, uh.AccessTokenDuration)
	resp := models.LoginResponse{
		AccessToken: token,
		User:        userResponse,
	}
	if err := respondWithJSON(w, http.StatusOK, resp); err != nil {
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
