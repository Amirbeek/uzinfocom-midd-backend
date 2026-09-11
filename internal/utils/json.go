package utils

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithPrivateFieldValidation())
}

func ReadJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024 * 4

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	dec.DisallowUnknownFields()
	return dec.Decode(data)
}

func WriteJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

type ErrorResponse struct {
	Error   string `json:"error"   example:"Bad Request"`
	Message string `json:"message" example:"invalid request body"`
}

func WriteJsonError(w http.ResponseWriter, status int, message string) error {
	return WriteJson(w, status, ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
