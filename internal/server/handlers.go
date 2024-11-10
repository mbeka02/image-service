package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/mbeka02/image-service/internal/database"
)

var validate *validator.Validate

// handle JSON responses
func ToJson(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(&payload)
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	params := CreateUserRequest{}
	json.NewDecoder(r.Body).Decode(&params)
	fmt.Println(params)
	validate := validator.New()
	if err := validate.Struct(params); err != nil {
		ToJson(w, http.StatusBadRequest, err.Error())
		return
	}
	passwordHash, err := HashPassword(params.Password)
	if err != nil {
		ToJson(w, http.StatusInternalServerError, err)
		return
	}
	user, err := s.Store.CreateUser(r.Context(), database.CreateUserParams{
		FullName: params.Fullname,
		Email:    params.Email,
		Password: passwordHash,
	})
	if err != nil {
		ToJson(w, http.StatusInternalServerError, err)
		return
	}
	ToJson(w, http.StatusCreated, user)
}
