package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// handle JSON responses
func ToJson(w http.ResponseWriter, statusCode int, payload interface{}) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	return json.NewEncoder(w).Encode(&payload)
}

func CreateUserAccount(w http.ResponseWriter, r *http.Request) {
	ToJson(w, 201, "account created")
}
