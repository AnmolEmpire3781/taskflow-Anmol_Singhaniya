package httpx

import (
    "encoding/json"
    "net/http"
)

type ErrorResponse struct {
    Error  string            `json:"error"`
    Fields map[string]string `json:"fields,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, status int, msg string) { JSON(w, status, ErrorResponse{Error: msg}) }
func ValidationError(w http.ResponseWriter, fields map[string]string) { JSON(w, http.StatusBadRequest, ErrorResponse{Error: "validation failed", Fields: fields}) }
